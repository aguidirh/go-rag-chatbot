package langchain

import (
	"log"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/tmc/langchaingo/llms/ollama"
)

type LanchChain struct {
	model          string
	scoreThreshold float32
	temperature    float64
	llm            *ollama.LLM
	kbCfg          data.KBConfig
}

func New(model string, scoreThreshold float32, temperature float64, kbCfg data.KBConfig) adapters.LLMHandler {

	llm, err := ollama.New(ollama.WithModel(model))
	if err != nil {
		log.Fatal(err)
	}

	return &LanchChain{model: model, scoreThreshold: scoreThreshold, temperature: temperature, llm: llm, kbCfg: kbCfg}
}
