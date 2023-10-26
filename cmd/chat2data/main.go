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
		docCMD(llm),
		cleanDocCacheCMD(),
		textCMD(llm),
		htmlCMD(llm),
		pdfCMD(llm),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("rootCmd err: %s", err)
	}
}

func runUI(qa qa.QA, info any) error {
	if cliMode {
		return cli.CLI(qa, info)
	}
	return web2.Web(webAddr, qa, info)
}
