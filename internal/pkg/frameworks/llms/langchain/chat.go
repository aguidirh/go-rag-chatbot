package langchain

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores"
)

func (l LanchChain) NewEmbedder() (emb embeddings.Embedder, err error) {
	e, err := embeddings.NewEmbedder(l.embeddingLlm)
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

	retriever := vectorstores.ToRetriever(vectorStore, 5, optionsVector...)
	// search
	resDocs, err := retriever.GetRelevantDocuments(context.Background(), query)

	stuffQAChain := chains.LoadStuffQA(l.chatLlm)
	answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": resDocs,
		"question":        query,
	}, options...)

	if err != nil {
		log.Fatalf("Error running chain: %v", err)
		return "", err
	}

	if answerText, ok := answer["text"]; ok {
		if text, isString := answerText.(string); isString {
			response = text
			for _, doc := range resDocs {
				if val, exists := doc.Metadata["url"]; exists {
					response += fmt.Sprintf("\n- Source: %s", val)
				}
			}
		} else {
			response = "unexpected error"
		}
	} else {
		log.Printf("Answer is not a string: %v", answer)
	}
	return response, nil
}
