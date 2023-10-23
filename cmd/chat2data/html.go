package main

import (
	"context"
	"fmt"
	"time"

	"github.com/byebyebruce/chat2data/htmlloader"
	"github.com/byebyebruce/chat2data/localvectordb"
	"github.com/byebyebruce/chat2data/qa/ragchain"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	embedding_openai "github.com/tmc/langchaingo/embeddings/openai"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

func htmlCMD(llm llms.LanguageModel) *cobra.Command {
	printDocsFlag := false
	cmd := &cobra.Command{
		Use:   "html url",
		Short: "html",
	}
	cmd.Flags().BoolVarP(&printDocsFlag, "print-docs", "p", false, "print docs")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Help()
		}
		url := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		docs, errLoad := htmlloader.Load(ctx, 1000, 0, url)
		if errLoad != nil {
			return errLoad
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

		const dbFile = "local_vector.db"
		db, err := localvectordb.New(dbFile)
		if err != nil {
			return err
		}
		defer db.Close()

		store, err := db.Get(e, url, 0)
		if err != nil {
			return err
		}
		oldDocs, err := store.Docs()
		if err != nil {
			return err
		}
		if len(oldDocs) != len(docs) {
			//refresh
			err := db.Delete(url)
			if err != nil {
				return err
			}
			store, err = db.Get(e, url, 0)
			if err != nil {
				return err
			}
			ctx1, cancel1 := context.WithTimeout(context.Background(), time.Minute*5)
			defer cancel1()

			err = store.AddDocuments(ctx1, docs)
			if err != nil {
				return err
			}
			fmt.Println("add docs", len(docs))
		} else {
			fmt.Println("use cache")
		}

		qa := ragchain.NewQAChain(llm, store, 5, 0.7)

		return runUI(qa, url)
	}

	return cmd
}

func printDocs(docs []schema.Document) {
	for i, doc := range docs {
		fmt.Println(color.RedString("Page #%d", i+1))
		fmt.Println(doc)
		fmt.Println()
		fmt.Println()
	}
}
