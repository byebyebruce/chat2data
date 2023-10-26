package main

import (
	"bytes"
	"path/filepath"

	"github.com/ledongthuc/pdf"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

func pdfCMD(llm llms.LanguageModel) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pdf <file>",
		Short: "pdf",
		RunE: docWrapper(llm, func(cmd *cobra.Command, args []string) (string, []schema.Document, error) {
			if len(args) != 1 {
				return "", nil, cmd.Help()
			}
			filePath, err := filepath.Abs(args[0])
			if err != nil {
				return "", nil, err
			}
			f, r, err := pdf.Open(filePath)
			// remember close file
			defer f.Close()
			if err != nil {
				return "", nil, err
			}
			var buf bytes.Buffer
			b, err := r.GetPlainText()
			if err != nil {
				return "", nil, err
			}
			buf.ReadFrom(b)

			docs, err := splitterText(filePath, buf.String())
			if err != nil {
				return "", nil, err
			}
			return filePath, docs, nil
		}),
	}
	docFlag(cmd)
	return cmd
}
