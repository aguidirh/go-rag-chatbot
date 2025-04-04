package qdrant

import (
	"fmt"
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

	store, err := NewQdrantStore(*url, emb, collectionName)
	if err != nil {
		return nil, fmt.Errorf("unable to create Qdrant store. %v", err)
	}

	return &QdrantDB{url: *url, store: *store}, nil
}

func NewQdrantStore(url url.URL, emb embeddings.Embedder, collectionName string) (*qdrant.Store, error) {
	store, err := qdrant.New(
		qdrant.WithURL(url),
		qdrant.WithCollectionName(collectionName),
		qdrant.WithEmbedder(emb),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to create Qdrant store. %v", err)
	}
	return &store, nil
}
