package docchain

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

func TestA(t *testing.T) {
	//f, err := os.Open("../testdata/sample.pdf")
	f, err := os.Open("../testdata/6.pdf")
	assert.NoError(t, err)
	defer f.Close()
	finfo, err := f.Stat()
	assert.NoError(t, err)
	p := documentloaders.NewPDF(f, finfo.Size())
	docs, err := p.Load(context.Background())
	var texts []string
	for _, doc := range docs {
		texts = append(texts, doc.PageContent)
	}

	splitter := textsplitter.NewRecursiveCharacter()
	splitter.ChunkSize = 1024
	splitter.ChunkOverlap = 0
	spdocs, err := textsplitter.CreateDocuments(splitter, texts, nil)
	assert.NoError(t, err)

	fmt.Println(spdocs)
}
