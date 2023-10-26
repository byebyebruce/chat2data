package htmlloader

import (
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/go-shiori/go-readability"
)

// ExtractContent extracts the content from a URL
func ExtractContent(timeout time.Duration, url string) (string, error) {
	article, err := readability.FromURL(url, timeout)
	if err != nil {
		return "", err
	}

	converter := md.NewConverter("", true, nil)
	in, err := converter.ConvertString(article.Content)
	if err != nil {
		return "", err
	}
	return in, nil
}
