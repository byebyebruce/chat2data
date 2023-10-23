package main

import (
	"flag"
	"log"

	"github.com/byebyebruce/chat2data/qa"
	"github.com/byebyebruce/chat2data/ui/cli"
	web2 "github.com/byebyebruce/chat2data/ui/web"
	"github.com/spf13/cobra"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
)

var (
	mysqlDSN    = flag.String("mysql", "", "mysql dsn (e.g. user:pwd@tcp(localhost:3306)/test)")
	sqlite3DSN  = flag.String("sqlite3", "", "sqlite3 dsn (e.g. test.db)")
	pgxDSN      = flag.String("postgre", "", "postgre dsn (e.g. postgres://db_user:mysecretpassword@localhost:5438/test?sslmode=disable)")
	csv         = flag.String("csv", "", "csv dir or file")
	useAllTable = flag.Bool("all", true, "use all table or choose by question")
	web         = flag.String("web", "", "web ui port")
)
var (
	webAddr string
	cliMode bool
)

func main() {
	flag.Parse()

	godotenv.Overload()

	llm, err := openai.NewChat()
	if err != nil {
		log.Fatalf("openai err: %s", err)
	}

	rootCmd := cobra.Command{
		Use:   "chat2data",
		Short: "chat2data",
	}
	rootCmd.PersistentFlags().StringVarP(&webAddr, "web", "w", "0.0.0.0:8088", "web ui listen address")
	rootCmd.PersistentFlags().BoolVarP(&cliMode, "cli", "c", false, "CLI mode")

	rootCmd.AddCommand(
		dbCMD(llm),
		csvCMD(llm),
		htmlCMD(llm),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("rootCmd err: %s", err)
	}

	/*
		var chain *dbchain.DBChain
		if *sqlite3DSN != "" {
			chain, err = dbchain.New(llm, sqlite3.EngineName, *sqlite3DSN, *useAllTable)
		} else if *mysqlDSN != "" {
			chain, err = dbchain.New(llm, mysql.EngineName, *mysqlDSN, *useAllTable)
		} else if *pgxDSN != "" {
			chain, err = dbchain.New(llm, postgresql.EngineName, *pgxDSN, *useAllTable)
		} else if *csv != "" {
			dbFile := path.Join(os.TempDir(), "chat2data.db")
			os.Remove(dbFile)
			defer os.Remove(dbFile)
			err = cmd.LoadCSV(dbFile, *csv)
			if err != nil {
				log.Fatalf("load csv err: %s", err)
			}
			chain, err = dbchain.New(llm, sqlite3.EngineName, dbFile, *useAllTable)
		} else {
			log.Fatalf("no dsn")
		}
		if err != nil {
			log.Fatalf("open database err: %s", err)
		}

		defer chain.Close()

		switch {
		case len(*web) > 0:
			web2.Web(*web, chain)
		default:
			if err := cmd.CLI(chain); err != nil {
				fmt.Println(err)
			}
		}

	*/
}

func runUI(qa qa.QA, info any) error {
	if cliMode {
		return cli.CLI(qa, info)
	}
	return web2.Web(webAddr, qa, info)
}
