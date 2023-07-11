package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/byebyebruce/chat2data/csv_importer"
	"github.com/byebyebruce/chat2data/util"
	"github.com/fatih/color"
	"github.com/sourcegraph/conc/pool"
	"github.com/tmc/langchaingo/tools/sqldatabase/sqlite3"
)

// LoadCSV load csv to sqlite3 database
func LoadCSV(dsn string, dirOrFile string) error {
	db, err := sql.Open(sqlite3.EngineName, dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	fi, err := os.Stat(dirOrFile)
	if err != nil {
		return err
	}
	var (
		fs []string
	)
	if fi.IsDir() {
		fs, err = util.FindFiles(dirOrFile, ".csv")
	} else {
		fs = []string{dirOrFile}
	}
	if err != nil {
		return err
	}
	if len(fs) == 0 {
		return fmt.Errorf("no csv file")
	}

	{
		bar := util.NewProgressBar(len(fs), "Insert table")
		p := pool.New().WithMaxGoroutines(runtime.NumCPU())
		mu := sync.Mutex{}
		errs := [][]string{}
		for i, f := range fs {
			_, f := i, f
			p.Go(func() {
				name, err := csv_importer.ImportCSV2DB(db, f)
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					errs = append(errs, []string{name, color.RedString(err.Error())})
					return
				}
				bar.Add(1)
			})
		}
		p.Wait()
		fmt.Println()
		if len(errs) > 0 {
			if len(errs) == len(fs) {
				return fmt.Errorf("all csv import failed")
			}
			util.RenderTable([]string{"table", "error"}, errs)
			fmt.Println()
			fmt.Println()
			return nil
		}
	}

	return nil
}
