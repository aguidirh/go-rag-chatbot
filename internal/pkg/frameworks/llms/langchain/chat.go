package langchain

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores"
)

func (l LanchChain) NewEmbedder() (emb embeddings.Embedder, err error) {
	e, err := embeddings.NewEmbedder(l.llm)
	if err != nil {
		log.Fatal(err) //TODO ALEX CHANGE ME
		return nil, err
	}

	return e, nil
}

func (l LanchChain) Chat(ctx context.Context, vectorStore vectorstores.VectorStore, query string) (response string, err error) {

	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(l.scoreThreshold),
	}

	options := []chains.ChainCallOption{
		chains.WithTemperature(l.temperature),
	}

	res, err := chains.Run(
		ctx,
		chains.NewRetrievalQAFromLLM(
			l.llm,
			vectorstores.ToRetriever(vectorStore, 10, optionsVector...),
		),
		query,
		options...,
	)

	if err != nil {
		log.Fatal(err) //TODO ALEX CHANGE ME
		return "", err
	}

	return res, nil
}
