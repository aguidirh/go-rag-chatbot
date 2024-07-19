package langchain

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

func (l LanchChain) DocumentLoader(ctx context.Context) ([]schema.Document, error) {

	var docs []schema.Document

	for _, doc := range l.kbCfg.Spec.Docs {
		ext := filepath.Ext(doc)
		switch ext {
		case ".txt", ".text":
		case ".md":
			if strings.Contains(doc, "http") {
				resp, err := httpCall(doc)
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()
				newDocs, err := documentloaders.NewText(resp.Body).LoadAndSplit(ctx, textsplitter.NewMarkdownTextSplitter())
				if err != nil {
					return nil, err
				}
				docs = append(docs, newDocs...)
			} else {
				//TODO ALEX implement file loader
			}
		case ".htm", ".html":
			if strings.Contains(doc, "http") {
				resp, err := httpCall(doc)
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()
				newDocs, err := documentloaders.NewHTML(resp.Body).LoadAndSplit(ctx, textsplitter.NewRecursiveCharacter())
				if err != nil {
					return nil, err
				}
				docs = append(docs, newDocs...)
			} else {
				//TODO ALEX implement file loader
			}

		}
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
