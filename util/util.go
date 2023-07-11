package util

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
)

func FindFiles(dir string, suffix string) ([]string, error) {
	var csvFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), suffix) {
			csvFiles = append(csvFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return csvFiles, err
}

func RenderTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetRowLine(true)
	table.SetCenterSeparator("*")
	table.SetAutoWrapText(false)
	//table.SetRowSeparator("-")

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func NewProgressBar(max int, desc string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(max,
		progressbar.OptionShowCount(),
		//progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		//progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan]"+desc+"..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	return bar
}
