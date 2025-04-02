package langchain

import (
	"context"
	"fmt"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
)

// DocumentLoader loads documents from the knowledge base configuration.
// It returns a slice of schema.Document and an error if any. If no collections are provided, it loads all documents.
func (l LanchChain) DocumentLoader(ctx context.Context, cb adapters.Crawlback, collections ...string) error {
	for _, doc := range l.kbCfg.Spec.Docs {
		if len(collections) > 0 {
			match := false
			for _, collection := range collections {
				if collection == doc.Collection {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		if len(doc.Type) == 0 {
			doc.Type = "http"
		}
		switch doc.Type {
		case "file":
			l.log.Warning("local files are not yet supported")
		case "http":
			for _, http := range doc.DocSourceHttp {
				err := l.crawler.Crawl(http, cb)
				if err != nil {
					return fmt.Errorf("unable to crawl http document of %s. %v", http.URL, err)
				}
			}
		}
	}

	return nil
}
