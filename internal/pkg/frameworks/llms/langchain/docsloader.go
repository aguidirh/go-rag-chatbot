package langchain

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/util"
	"github.com/stretchr/objx"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

// LoadDocumentsFromText loads documents from a given text string and metadata.
func (l LanchChain) LoadDocumentsFromHttpRequest(ctx context.Context, cb adapters.Crawlback, text string, r *http.Request) error {
	docs, err := documentloaders.NewText(r.Body).LoadAndSplit(ctx, textsplitter.NewRecursiveCharacter())

	if err != nil {
		return fmt.Errorf("unable to load and split the text part of the document. %v", err)
	}

	metadataParam := util.GetQueryParameterAsString(r, "metadata", "{}")
	metadata, err := objx.FromJSON(metadataParam)
	if err != nil {
		return fmt.Errorf("unable to parse metadata. %v", err)
	}

	for _, doc := range docs {
		for k, v := range metadata {
			doc.Metadata[k] = v
		}
	}
	err = cb(docs, nil)
	if err != nil {
		return fmt.Errorf("error encountered in callback function. %v", err)
	}
	return err
}

// LoadDocumentsFromConfig loads documents from the knowledge base configuration.
// It returns a slice of schema.Document and an error if any. If no collections are provided, it loads all documents.
func (l LanchChain) LoadDocumentsFromConfig(ctx context.Context, cb adapters.Crawlback, collections ...string) error {
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
