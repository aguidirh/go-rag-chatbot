package adapters

import (
	"context"
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/tmc/langchaingo/embeddings"   //TODO it should not be specific to langchain
	"github.com/tmc/langchaingo/schema"       //TODO it should not be specific to langchain
	"github.com/tmc/langchaingo/vectorstores" //TODO it should not be specific to langchain
)

type VectorDB interface {
	CreateCollection(ctx context.Context, collectionName string, vectorSize int, distance string) error
	DeleteCollection(ctx context.Context, collectionName string) error
	DoesCollectionExist(ctx context.Context, collectionName string) (bool, error)
	AddDocuments(ctx context.Context, docs []schema.Document) error
	GetStore() vectorstores.VectorStore
}

type Crawlback func(docs []schema.Document, e *colly.HTMLElement) error

type LLMHandler interface {
	NewEmbedder() (emb embeddings.Embedder, err error)
	LoadDocumentsFromConfig(ctx context.Context, cb Crawlback, collections ...string) error
	LoadDocumentsFromHttpRequest(ctx context.Context, cb Crawlback, collection string, r *http.Request) error
	Chat(ctx context.Context, vectorStore vectorstores.VectorStore, query string) (response string, err error)
}

type App interface {
	Chat(ctx context.Context, askMeSomething string) (string, error)
}
