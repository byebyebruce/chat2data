package ragchain

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/byebyebruce/chat2data/localvectordb"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

type RAGChain struct {
	llm       llms.LanguageModel
	retriever schema.Retriever
}

func NewQAChain(llm llms.LanguageModel, store vectorstores.VectorStore, scoreTopN int, scoreThreshold float64) *RAGChain {
	retriever := vectorstores.ToRetriever(
		store,
		scoreTopN,
		vectorstores.WithScoreThreshold(scoreThreshold))
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

func RefreshDoc(db localvectordb.DB, embedder embeddings.Embedder, name string, docs []schema.Document) (bool, error) {
	store, err := db.Get(embedder, name, 0)
	if err != nil {
		return false, err
	}
	oldDocs, err := store.Docs()
	if err != nil {
		return false, err
	}
	needRefresh := false
	if len(oldDocs) != len(docs) {
		needRefresh = true
	} else {
		for i, doc := range oldDocs {
			oldmd5 := md5.Sum([]byte(doc.PageContent))
			newmd5 := md5.Sum([]byte(docs[i].PageContent))
			if !bytes.Equal(oldmd5[:], newmd5[:]) {
				needRefresh = true
				break
			}
		}
	}

	if needRefresh {
		err := db.Delete(name)
		if err != nil {
			return false, err
		}
		store, err = db.Get(embedder, name, 0)
		if err != nil {
			return false, err
		}
		ctx1, cancel1 := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel1()

		for i := range docs {
			id := fmt.Sprintf("%s|%04d", name, i) // 为了得到稳定的排序
			localvectordb.PutMetaID(&docs[i].Metadata, id)
		}
		err = store.AddDocuments(ctx1, docs)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

type DocRAGChain struct {
	*RAGChain
}

func NewDocRAGChain(llm llms.LanguageModel, db localvectordb.DB, embedder embeddings.Embedder, name string, scoreTopN int, scoreThreshold float64) (*DocRAGChain, error) {
	store, err := db.Get(embedder, name, 0)
	if err != nil {
		return nil, err
	}
	return &DocRAGChain{
		RAGChain: NewQAChain(llm, store, scoreTopN, scoreThreshold),
	}, nil
}
