package langchain

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

func (l LanchChain) getHttpDocument(ctx context.Context, url string) ([]schema.Document, error) {
	resp, err := httpCall(url)
	if err != nil {
		return nil, fmt.Errorf("unable to get document(s) from %s. %v", url, err)
	}
	defer resp.Body.Close()
	newDocs, err := documentloaders.NewText(resp.Body).LoadAndSplit(ctx, textsplitter.NewMarkdownTextSplitter())
	if err != nil {
		return nil, fmt.Errorf("unable to load document(s) from %s. %v", url, err)
	}
	return newDocs, nil
}

func (l LanchChain) DocumentLoader(ctx context.Context) ([]schema.Document, error) {

	var docs []schema.Document
	var err error
	for _, doc := range l.kbCfg.Spec.Docs {
		if len(doc.Type) == 0 {
			doc.Type = "http"
		}
		switch doc.Type {
		case "file":
			l.log.Warning("local files are not yet supported")
		case "http":
			for _, http := range doc.DocSourceHttp {
				if http.RecursionLevels > 0 {
					descendants, err := l.crawler.GetDirectDescendants(http.URL, http.RecursionLevels)
					if err != nil {
						return nil, fmt.Errorf("unable to get descendants of %s. %v", http.URL, err)
					}
					for descendant := range descendants {
						descendantDocs, err := l.getHttpDocument(ctx, descendant)
						if err != nil {
							return nil, fmt.Errorf("unable to get http document of %s. %v", http.URL, err)
						}
						docs = append(docs, descendantDocs...)
					}
				} else {
					docs, err = l.getHttpDocument(ctx, http.URL)
					if err != nil {
						return nil, fmt.Errorf("unable to get http document of %s. %v", http.URL, err)
					}
				}
			}
		}
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

	return docs, nil

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
}

func httpCall(resource string) (*http.Response, error) {
	resp, err := http.Get(resource)

	if err != nil {
		return resp, err
	}

	return resp, nil
}
