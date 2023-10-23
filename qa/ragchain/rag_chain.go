package ragchain

import (
	"context"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

type RAGChain struct {
	llm       llms.LanguageModel
	retriever schema.Retriever
}

func NewQAChain(llm llms.LanguageModel, store vectorstores.VectorStore, docsNumber int, scoreThreshold float64) *RAGChain {
	retriever := vectorstores.ToRetriever(store, docsNumber, vectorstores.WithScoreThreshold(scoreThreshold))
	return &RAGChain{
		llm:       llm,
		retriever: retriever,
	}
}

func (q *RAGChain) Answer(ctx context.Context, question string) (string, error) {
	answer, err := chains.Run(ctx,
		chains.NewRetrievalQAFromLLM(q.llm, q.retriever),
		question)
	return answer, err
}
