package qdrant

import (
	"fmt"
	"log"
	"net/url"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

type QdrantDB struct {
	url   url.URL
	store qdrant.Store
}

func New(host, port string, collectionName string, emb embeddings.Embedder) (adapters.VectorDB, error) {
	url, err := url.Parse(fmt.Sprintf("http://%s:%s", host, port))
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL. %v", err)
	}

	store := NewQdrantStore(*url, emb, collectionName)

	return &QdrantDB{url: *url, store: store}, nil
}

func NewQdrantStore(url url.URL, emb embeddings.Embedder, collectionName string) qdrant.Store {
	store, err := qdrant.New(
		qdrant.WithURL(url),
		qdrant.WithCollectionName(collectionName),
		qdrant.WithEmbedder(emb),
	)

	if err != nil {
		log.Fatal(err) //TODO ALEX CHANGES ME
	}

	return store
}
