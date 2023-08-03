//go:generate protoc --proto_path=./ --gogo_out=paths=source_relative:./ types.proto
package localvectordb

import (
	"container/heap"
	"context"
	"encoding/json"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/sourcegraph/conc/pool"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

var _ vectorstores.VectorStore = (*Store)(nil)

type Store struct {
	embedder  embeddings.Embedder
	db        *bolt.DB
	maxGo     int
	nameSpace string
}

func New(path string, embedder embeddings.Embedder, namespace string, maxGo int) (Store, error) {
	db, err := bolt.Open(path, os.ModePerm, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return Store{}, err
	}
	if maxGo == 0 {
		maxGo = runtime.GOMAXPROCS(0)
	}

	s := Store{
		embedder:  embedder,
		db:        db,
		maxGo:     maxGo,
		nameSpace: namespace,
	}

	if err := s.CreateNamespace(s.nameSpace); err != nil {
		return s, nil
	}
	return s, nil
}

func (s Store) AddDocuments(ctx context.Context, documents []schema.Document, option ...vectorstores.Option) error {
	opt := s.getOption(option...)
	texts := make([]string, 0, len(documents))
	for _, document := range documents {
		texts = append(texts, document.PageContent)
	}
	vectors, err := s.embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return err
	}

	docs := make([]*Doc, 0, len(documents))
	for i, document := range documents {
		doc := &Doc{
			Id:      uuid.New().String(),
			Content: document.PageContent,
			Vector:  convertVector(vectors[i]),
		}
		if document.Metadata != nil {
			doc.Meta, _ = json.Marshal(document.Metadata)
		}
		docs = append(docs, doc)
	}

	bucket := []byte(opt.NameSpace)
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		for i := 0; i < len(docs); i++ {
			bs, err := proto.Marshal(docs[i])
			if err != nil {
				return err
			}
			if err = b.Put([]byte(docs[i].Id), bs); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s Store) getOption(options ...vectorstores.Option) vectorstores.Options {
	opt := vectorstores.Options{
		NameSpace:      s.nameSpace,
		ScoreThreshold: 0,
	}
	for _, option := range options {
		option(&opt)
	}
	return opt
}

func (s Store) SimilaritySearch(ctx context.Context, query string, numDocuments int, options ...vectorstores.Option) ([]schema.Document, error) {
	opt := s.getOption(options...)
	vector64, err := s.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	var (
		vector = convertVector(vector64)
		bucket = []byte(opt.NameSpace)
		h      = &bigHeap{}
		p      = pool.New().WithMaxGoroutines(s.maxGo)
		mu     = &sync.Mutex{}
	)

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.ForEach(func(k, v []byte) error {
			p.Go(func() {
				d := &Doc{}
				if err := proto.Unmarshal(v, d); err != nil {
					return
				}
				score := Cosine(vector, d.Vector)
				if float64(score) < opt.ScoreThreshold {
					return
				}
				mu.Lock()
				heap.Push(h, &DocWithScore{
					Doc:   d,
					Score: score,
				})
				if len(*h) > numDocuments {
					heap.Pop(h)
				}
				mu.Unlock()
			})
			return nil
		})
	})
	p.Wait()

	n := len(*h)
	if n == 0 {
		return nil, err
	}
	docs := make([]schema.Document, n)
	for i := n - 1; i >= 0; i-- {
		ds := heap.Pop(h).(*DocWithScore)
		d := schema.Document{
			PageContent: ds.Content,
		}
		if len(ds.Meta) > 0 {
			json.Unmarshal(ds.Meta, &d.Metadata)
		}
		docs[i] = d
	}
	return docs, nil
}

func (s Store) ListNamespace() ([]string, error) {
	var ns []string
	err := s.db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			ns = append(ns, string(name))
			return nil
		})
		return nil
	})
	return ns, err
}

func (s Store) CreateNamespace(namespace string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(namespace))
		return err
	})
}

func (s Store) DeleteNamespace(namespace string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(namespace))
	})
}

func (s Store) Close() error {
	return s.db.Close()
}

func convertVector(v []float64) []float32 {
	v1 := make([]float32, len(v))
	for i := 0; i < len(v); i++ {
		v1[i] = float32(v[i])
	}
	return v1
}
