package util

import (
	"fmt"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/databases/qdrant"
	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/tmc/langchaingo/embeddings"
)

// GetVectorDBForCollection creates a new VectorDB client for a specified collection
func GetVectorDBForCollection(collectionName string, vectorDBConfig *data.VectorDB, embedder embeddings.Embedder) (adapters.VectorDB, error) {
	// Create a new Qdrant client using the provided configuration and embedder
	client, err := qdrant.New(vectorDBConfig.Host, vectorDBConfig.Port, collectionName, embedder)
	if err != nil {
		return nil, fmt.Errorf("failed to create VectorDB client: %w", err)
	}
	return client, nil
}
