package langchain

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores"
)

func (l LanchChain) NewEmbedder() (emb embeddings.Embedder, err error) {

	if l.embeddingLlm != nil {
		emb, err = embeddings.NewEmbedder(l.embeddingLlm)
		if err != nil {
			return nil, fmt.Errorf("failed to create embedding model: %v", err)
		}
	} else if l.openAIEmbeddingLLM != nil {
		emb, err = embeddings.NewEmbedder(l.openAIEmbeddingLLM)
		if err != nil {
			return nil, fmt.Errorf("failed to create embedding model: %v", err)
		}
	}
	if emb == nil {
		return nil, errors.New("no known embedding model provided")
	}
	return emb, nil
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
	var stuffQAChain chains.StuffDocuments
	if l.chatLlm != nil {
		stuffQAChain = chains.LoadStuffQA(l.chatLlm)
	} else if l.openAIChatLLM != nil {
		stuffQAChain = chains.LoadStuffQA(l.openAIChatLLM)
	}

	answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": resDocs,
		"question":        query,
	}, options...)

	if err != nil {
		return "", fmt.Errorf("Error running chain: %v", err)
	}

	if answerText, ok := answer["text"]; ok {
		if text, isString := answerText.(string); isString {
			response = text
			for _, doc := range resDocs {
				if val, exists := doc.Metadata["url"]; exists {
					response += fmt.Sprintf("\n- source: %s", val)
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
