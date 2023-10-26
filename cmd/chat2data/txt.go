package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

func textCMD(llm llms.LanguageModel) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txt <file>",
		Short: "txt",
		RunE: docWrapper(llm, func(cmd *cobra.Command, args []string) (string, []schema.Document, error) {
			if len(args) != 1 {
				return "", nil, cmd.Help()
			}
			filePath, err := filepath.Abs(args[0])
			if err != nil {
				return "", nil, err
			}

			b, err := os.ReadFile(filePath)
			if err != nil {
				return "", nil, err
			}

			docs, err := splitterText(filePath, string(b))
			if err != nil {
				return "", nil, err
			}

			return filePath, docs, nil
		}),
	}
	docFlag(cmd)
	return cmd
}
