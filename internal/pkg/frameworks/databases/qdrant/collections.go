package qdrant

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/stretchr/objx"
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
		return fmt.Errorf("unable to create collection. %v", err)
	}

	return nil
}

func (q QdrantDB) DoesCollectionExist(ctx context.Context, collectionName string) (bool, error) {
	collectionConfig := map[string]interface{}{}
	reader, _, err := qdrant.DoRequest(ctx, *q.url.JoinPath(qdrant_create_collection_api, collectionName, "exists"), "", http.MethodGet, collectionConfig)
	if err != nil {
		return false, fmt.Errorf("unable to check if collection exists. %v", err)
	}
	defer reader.Close()
	content, err := io.ReadAll(reader)
	if err != nil {
		return false, fmt.Errorf("unable to read response body. %v", err)
	}

	m, err := objx.FromJSON(string(content))
	if err != nil {
		return false, fmt.Errorf("unable to parse JSON response body. %v", err)
	}

	return m.Get("result.exists").Bool(false), nil
}

func (q QdrantDB) DeleteCollection(ctx context.Context, collectionName string) error {

	exists, err := q.DoesCollectionExist(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("unable to check if collection exists. %v", err)
	}

	if !exists {
		return fmt.Errorf("collection does not exist. %v", err)
	}

	collectionConfig := map[string]interface{}{}
	_, _, err = qdrant.DoRequest(ctx, *q.url.JoinPath(qdrant_create_collection_api, collectionName), "", http.MethodDelete, collectionConfig)

	if err != nil {
		return fmt.Errorf("unable to delete collection. %v", err)
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
