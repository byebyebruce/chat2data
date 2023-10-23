package htmlloader

import (
	"context"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/go-shiori/go-readability"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

// Load loads documents from urls
func Load(ctx context.Context, chunkSize, chunkOverlap int, urls ...string) ([]schema.Document, error) {
	docs := []schema.Document{}
	for i := 0; i < len(urls); i++ {
		url := urls[i]

		article, err := readability.FromURL(url, 30*time.Second)
		if err != nil {
			return nil, err
		}

		converter := md.NewConverter("", true, nil)
		in, err := converter.ConvertString(article.Content)
		if err != nil {
			return nil, err
		}

		ts := textsplitter.NewRecursiveCharacter()
		ts.ChunkOverlap = chunkOverlap
		ts.ChunkSize = chunkSize
		texts, err := ts.SplitText(in)
		if err != nil {
			return nil, err
		}
		for _, text := range texts {
			if strings.TrimSpace(text) == "" {
				continue
			}
			docs = append(docs,
				schema.Document{
					PageContent: text,
					Metadata: map[string]interface{}{
						"url": url,
						"idx": i,
					},
				})
		}
	}
	return docs, nil
}
