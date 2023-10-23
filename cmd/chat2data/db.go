package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/byebyebruce/chat2data/qa"
	"github.com/byebyebruce/chat2data/qa/dbchain"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools/sqldatabase/mysql"
	"github.com/tmc/langchaingo/tools/sqldatabase/postgresql"
	"github.com/tmc/langchaingo/tools/sqldatabase/sqlite3"
)

func dbCMD(llm llms.LanguageModel) *cobra.Command {
	useAllTables := false
	cmd := &cobra.Command{
		Use:   "db dsn",
		Short: "db",
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
		dsn := args[0]
		dataSource := ""
		if strings.HasPrefix(dsn, "postgres://") {
			qa, err = dbchain.New(llm, postgresql.EngineName, dsn, useAllTables)
			dataSource = fmt.Sprintf("postgre: %s", dsn)
		} else {
			if s, err := os.Stat(dsn); err == nil && !s.IsDir() {
				qa, err = dbchain.New(llm, sqlite3.EngineName, dsn, useAllTables)
				dataSource = fmt.Sprintf("sqlite: %s", dsn)
			} else {
				qa, err = dbchain.New(llm, mysql.EngineName, dsn, useAllTables)
				dataSource = fmt.Sprintf("mysql: %s", dsn)
			}
		}
		if err != nil {
			return err
		}
		return runUI(qa, dataSource)
	}

	return cmd
}
