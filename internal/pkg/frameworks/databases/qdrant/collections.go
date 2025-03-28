package qdrant

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

func (q QdrantDB) CreateCollection(ctx context.Context, collectionName string, vectorSize int, distance string) error {

	collectionConfig := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     vectorSize,
			"distance": distance,
		},
	}

	_, _, err := qdrant.DoRequest(ctx, *q.url.JoinPath(qdrant_create_collection_api, collectionName), "", http.MethodPut, collectionConfig)

	if err != nil {
		return err
	}

	return nil
}

func (q QdrantDB) AddDocuments(ctx context.Context, docs []schema.Document) error {

	var wg sync.WaitGroup

	if len(docs) > 0 {
		for _, doc := range docs {

			wg.Add(1)
			go func(doc schema.Document) { //TODO add the done channel to avoid goroutine leak
				defer wg.Done()
				//vectorstores.Options.Embedder.EmbedDocuments()
				_, err := q.store.AddDocuments(ctx, docs)
				if err != nil {
					log.Fatal(err) //TODO ALEX CHANGES ME
				}
			}(doc)
		}
	}

	wg.Wait()
	return nil
}

func (q QdrantDB) GetStore() vectorstores.VectorStore {
	return q.store
}
