package main

import (
	"time"

	"github.com/byebyebruce/chat2data/htmlloader"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

func htmlCMD(llm llms.LanguageModel) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "html <url>",
		Short: "html",
		RunE: docWrapper(llm, func(cmd *cobra.Command, args []string) (string, []schema.Document, error) {
			if len(args) != 1 {
				return "", nil, cmd.Help()
			}
			url := args[0]

			c, err := htmlloader.ExtractContent(time.Second*30, url)
			if err != nil {
				return "", nil, err
			}
			docs, err := splitterText(url, c)
			if err != nil {
				return "", nil, err
			}
			return url, docs, nil
		}),
	}
	docFlag(cmd)
	return cmd
}
