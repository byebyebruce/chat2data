package main

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/byebyebruce/chat2data/csvimporter"
	"github.com/byebyebruce/chat2data/qa"
	"github.com/byebyebruce/chat2data/qa/dbchain"
	"github.com/byebyebruce/chat2data/util"
	"github.com/fatih/color"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools/sqldatabase/sqlite3"
)

func csvCMD(llm llms.LanguageModel) *cobra.Command {
	useAllTables := false
	cmd := &cobra.Command{
		Use:   "csv file|dir",
		Short: "csv",
	}
	cmd.Flags().BoolVarP(&useAllTables, "all-tables", "a", false, "use all tables")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Help()
		}

		var (
			err error
			qa  qa.QA
		)
		fileOrDir := args[0]
		const tmpSQLiteFile = "./chat.db"
		err = csv2SQLite(fileOrDir, tmpSQLiteFile)
		if err != nil {
			return err
		}
		defer os.Remove(tmpSQLiteFile)

		qa, err = dbchain.New(llm, sqlite3.EngineName, tmpSQLiteFile, useAllTables)
		if err != nil {
			return err
		}
		dataSource := fmt.Sprintf("csv: %s", fileOrDir)
		return runUI(qa, dataSource)
	}

	return cmd
}

// csv2SQLite load csv to sqlite3 database
func csv2SQLite(dirOrFile string, outSQLiteFile string) error {
	db, err := sql.Open(sqlite3.EngineName, outSQLiteFile)
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
				name, err := csvimporter.ImportCSV2DB(db, f)
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
