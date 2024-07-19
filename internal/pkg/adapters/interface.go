package adapters

import (
	"context"

	"github.com/tmc/langchaingo/embeddings"   //TODO it should not be specific to langchain
	"github.com/tmc/langchaingo/schema"       //TODO it should not be specific to langchain
	"github.com/tmc/langchaingo/vectorstores" //TODO it should not be specific to langchain
)

type VectorDB interface {
	CreateCollection(ctx context.Context, collectionName string, vectorSize int, distance string) error
	AddDocuments(ctx context.Context, docs []schema.Document) error
	GetStore() vectorstores.VectorStore
}

type LLMHandler interface {
	NewEmbedder() (emb embeddings.Embedder, err error)
	DocumentLoader(ctx context.Context) (docs []schema.Document, err error)
	Chat(ctx context.Context, vectorStore vectorstores.VectorStore, query string) (response string, err error)
}

type App interface {
	Chat(ctx context.Context, askMeSomething string) (string, error)
}
