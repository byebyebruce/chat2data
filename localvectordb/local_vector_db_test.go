package localvectordb

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/embeddings/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

func TestStore_SimilaritySearch(t *testing.T) {
	//cli, err := openai.New(openai.WithModel(""))
	//require.NoError(t, err)
	e, err := openai.NewOpenAI()

	require.NoError(t, err)

	const dbFile = "a.db"
	defer os.Remove(dbFile)
	db, err := New(dbFile)
	defer db.Close()
	require.NoError(t, err)

	store, err := db.Get(e, "test", 0)
	require.NoError(t, err)
	err = store.AddDocuments(context.Background(), []schema.Document{
		{PageContent: "Tokyo"},
		{PageContent: "Yokohama"},
		{PageContent: "Osaka"},
		{PageContent: "Nagoya"},
		{PageContent: "Sapporo"},
		{PageContent: "Fukuoka"},
		{PageContent: "Dublin"},
		{PageContent: "Paris"},
		{PageContent: "London "},
		{PageContent: "New York"},
	})
	require.NoError(t, err)

	// test with a score threshold of 0.8, expected 6 documents
	docs, err := store.SimilaritySearch(context.Background(),
		"Which of these are cities in Japan", 10,
		vectorstores.WithScoreThreshold(0.8))
	require.NoError(t, err)
	//require.Len(t, docs, 6)
	require.Equal(t, "Tokyo", docs[0].PageContent)

	// test with a score threshold of 0, expected all 10 documents
	docs, err = store.SimilaritySearch(context.Background(),
		"Which of these are cities in Japan", 10,
		vectorstores.WithScoreThreshold(0))
	require.NoError(t, err)
	require.Len(t, docs, 10)
}
