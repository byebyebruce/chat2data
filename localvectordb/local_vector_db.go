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

type DB struct {
	db *bolt.DB
}

func New(path string) (DB, error) {
	db, err := bolt.Open(path, os.ModePerm, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return DB{}, err
	}

	s := DB{
		db: db,
	}

	return s, nil
}

func (s DB) Get(embedder embeddings.Embedder, namespace string, maxGo int) (Store, error) {
	if maxGo == 0 {
		maxGo = runtime.GOMAXPROCS(0)
	}

	store := Store{
		embedder: embedder,
		db:       s.db,
		maxGo:    maxGo,
		name:     namespace,
	}

	if err := s.create(namespace); err != nil {
		return store, err
	}
	return store, nil
}

func (s DB) List() ([]string, error) {
	var ns []string
	err := s.db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			ns = append(ns, string(name))
			return nil
		})
		for i := len(ns) - 1; i >= 0; i-- {
			n := ns[i]
			c := tx.Bucket([]byte(n)).Cursor()
			if k, _ := c.First(); k == nil {
				ns = append(ns[:i], ns[i+1:]...)
			}
		}
		return nil
	})
	return ns, err
}

func (s DB) create(namespace string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(namespace))
		return err
	})
}

func (s DB) Delete(namespace string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(namespace))
	})
}

func (s DB) Close() error {
	return s.db.Close()
}

func convertVector(v []float64) []float32 {
	v1 := make([]float32, len(v))
	for i := 0; i < len(v); i++ {
		v1[i] = float32(v[i])
	}
	return v1
}

var _ vectorstores.VectorStore = (*Store)(nil)

type Store struct {
	name     string
	embedder embeddings.Embedder
	db       *bolt.DB
	maxGo    int
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
			id, ok := GetMetaID(document.Metadata)
			if ok {
				// 为了得到稳定的排序
				doc.Id = id
			}
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
		NameSpace:      s.name,
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

func (s Store) Docs() ([]schema.Document, error) {
	var ns []schema.Document
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.name))
		b.ForEach(func(k []byte, v []byte) error {
			d := Doc{}
			if err := proto.Unmarshal(v, &d); err != nil {
				return err
			}
			var m map[string]interface{}
			if len(d.Meta) > 0 {
				json.Unmarshal(d.Meta, &m)
			}
			ns = append(ns, schema.Document{PageContent: d.Content, Metadata: m})
			return nil
		})

		return nil
	})
	return ns, err
}

func PutMetaID(m *map[string]any, id string) {
	if *m == nil {
		*m = make(map[string]any)
	}
	(*m)["_id"] = id
}

func GetMetaID(m map[string]any) (string, bool) {
	if m == nil {
		return "", false
	}
	if a, ok := m["_id"]; !ok {
		return "", false
	} else {
		id, ok := a.(string)
		if !ok {
			return "", false
		}
		return id, true
	}
}
