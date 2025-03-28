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
				err := l.crawler.Crawl(http.URL, http.RecursionLevels, cb, http.AllowedDomains)
				if err != nil {
					return fmt.Errorf("unable to crawl http document of %s. %v", http.URL, err)
				}
			}
		}
	}

	return nil

	// f, err := os.Open("/home/aguidi/go/src/github.com/aguidirh/go-rag-chatbot/test.txt")
	// if err != nil {
	// 	fmt.Println("Error opening file: ", err)
	// }

	// docs, err = documentloaders.NewText(f).LoadAndSplit(ctx, textsplitter.NewRecursiveCharacter())
	// // fmt.Println(docs_text)

	// if err != nil {
	// 	log.Fatal(err) //TODO ALEX CHANGE ME
	// 	return nil, err
	// }

	// case ".htm", ".html":
	// 	if strings.Contains(doc, "http") {
	// 		resp, err := httpCall(doc)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		defer resp.Body.Close()
	// 		newDocs, err := documentloaders.NewHTML(resp.Body).LoadAndSplit(ctx, textsplitter.NewRecursiveCharacter())
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		docs = append(docs, newDocs...)
	// 	} else {
	// 		//TODO ALEX implement file loader
	// 	}

	// }

}
