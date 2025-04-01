package langchain

import (
	"log"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/corpus/crawler"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/util"
	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/llms/ollama"
)

type LanchChain struct {
	model          string
	embeddingModel string
	scoreThreshold float32
	temperature    float64
	llm            *ollama.LLM
	embeddingsLlm  *ollama.LLM
	kbCfg          data.KBConfig
	log            *logrus.Logger
	crawler        *crawler.Crawler
	http           *util.HttpAccessor
}

func New(model, embeddingModel, url string,
	scoreThreshold float32,
	temperature float64,
	kbCfg data.KBConfig,
	logger *logrus.Logger) adapters.LLMHandler {
	if len(url) == 0 {
		url = "http://127.0.0.1:11434"
	}
	llm, err := ollama.New(ollama.WithModel(model), ollama.WithServerURL(url))
	if err != nil {
		log.Fatal(err)
	}
	embeddingsLlm, err := ollama.New(ollama.WithModel(embeddingModel), ollama.WithServerURL(url))
	if err != nil {
		log.Fatal(err)
	}

	return &LanchChain{model: model, embeddingModel: embeddingModel, scoreThreshold: scoreThreshold, temperature: temperature, llm: llm, embeddingsLlm: embeddingsLlm, kbCfg: kbCfg, log: logger, crawler: crawler.New(logger), http: util.NewHttpAccessor()}
}
