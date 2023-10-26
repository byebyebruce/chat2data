package main

import (
	"fmt"
	"os"
	"path"

	"github.com/byebyebruce/chat2data/localvectordb"
	"github.com/byebyebruce/chat2data/qa/ragchain"
	"github.com/byebyebruce/chat2data/ui/cli"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	embedding_openai "github.com/tmc/langchaingo/embeddings/openai"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

const dbFile = ".chat2data.tmp"

var (
	tmpDBFilePath string
)
var (
	printDocsFlag  bool
	topN           int
	scoreThreshold float64
	chunkSize      int
	chunkOverlap   int
)

func init() {
	dir, err := os.UserHomeDir()
	if err != nil {
		dir = os.TempDir()
	}
	tmpDBFilePath = path.Join(dir, dbFile)
}

func docFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVarP(&chunkSize, "chunk-size", "s", 1000, "chunk size")
	cmd.PersistentFlags().IntVarP(&chunkOverlap, "chunk-overlap", "o", 0, "chunk overlap")
	cmd.PersistentFlags().IntVarP(&topN, "topN", "n", 5, "vector search topN")
	cmd.PersistentFlags().BoolVarP(&printDocsFlag, "print-docs", "p", false, "print docs")
	cmd.PersistentFlags().Float64VarP(&scoreThreshold, "score-threshold", "t", 0.7, "score threshold")
}

func printDocs(docs []schema.Document) {
	for i, doc := range docs {
		fmt.Println(color.RedString("Page #%d", i+1))
		fmt.Println(doc)
		fmt.Println()
		fmt.Println()
	}
}

func splitterText(from string, str string) ([]schema.Document, error) {
	ts := textsplitter.NewRecursiveCharacter()
	ts.ChunkOverlap = chunkOverlap
	ts.ChunkSize = chunkSize
	chunks, err := ts.SplitText(str)
	if err != nil {
		return nil, err
	}
	var docs []schema.Document
	for _, c := range chunks {
		docs = append(docs, schema.Document{
			PageContent: c,
		})
	}
	return docs, nil
}

type RunE func(cmd *cobra.Command, args []string) error

func docWrapper(llm llms.LanguageModel, f func(cmd *cobra.Command, args []string) (string, []schema.Document, error)) RunE {
	return func(cmd *cobra.Command, args []string) error {
		name, docs, err := f(cmd, args)
		if err != nil {
			return err
		}
		if len(docs) == 0 {
			return fmt.Errorf("no docs")
		}
		if printDocsFlag {
			printDocs(docs)
		}

		e, err := embedding_openai.NewOpenAI()
		if err != nil {
			return err
		}

		db, err := localvectordb.New(tmpDBFilePath)
		if err != nil {
			return err
		}
		fmt.Println("load cache db file", tmpDBFilePath)
		defer db.Close()

		ok, err := ragchain.RefreshDoc(db, e, name, docs)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("use cached doc")
		}
		qa, err := ragchain.NewDocRAGChain(llm, db, e, name, topN, scoreThreshold)
		if err != nil {
			return err
		}

		return runUI(qa, name)
	}
}

func docCMD(llm llms.LanguageModel) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "Choose doc to chat(need add doc to cache first)",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		e, err := embedding_openai.NewOpenAI()
		if err != nil {
			return err
		}

		db, err := localvectordb.New(tmpDBFilePath)
		if err != nil {
			return err
		}
		defer db.Close()
		names, err := db.List()
		if err != nil {
			return err
		}
		if len(names) == 0 {
			fmt.Println("no doc in cache")
			return nil
		}
		for {
			sel := promptui.Select{
				Label: "Select doc",
				Items: names,
			}
			_, name, err := sel.Run()
			if err != nil {
				if err == promptui.ErrInterrupt {
					return nil
				}
				return err
			}

			qa, err := ragchain.NewDocRAGChain(llm, db, e, name, topN, scoreThreshold)
			if err != nil {
				return err
			}

			err = cli.CLI(qa, name)
			if err != nil {
				return err
			}
		}
	}
	return cmd
}

func cleanDocCacheCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "clean doc cache",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := os.RemoveAll(tmpDBFilePath)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("cache cleaned", tmpDBFilePath)
		}
		return nil
	}
	return cmd
}
