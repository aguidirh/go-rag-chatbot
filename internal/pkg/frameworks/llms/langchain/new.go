package langchain

import (
	"log"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/corpus/crawler"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/util"
	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

type LanchChain struct {
	embeddingModel     string
	chatModel          string
	scoreThreshold     float32
	temperature        float64
	embeddingLlm       *ollama.LLM
	chatLlm            *ollama.LLM
	openAIEmbeddingLLM *openai.LLM
	openAIChatLLM      *openai.LLM
	kbCfg              data.KBConfig
	log                *logrus.Logger
	crawler            *crawler.Crawler
	http               *util.HttpAccessor
}

func New(chatModel, embeddingModel, url string,
	scoreThreshold float32,
	temperature float64,
	kbCfg data.KBConfig,
	cfg data.Config,
	logger *logrus.Logger) adapters.LLMHandler {
	if len(url) == 0 {
		url = "http://127.0.0.1:11434"
	}

	switch cfg.Spec.LLM.ProviderType {
	case data.LLMProviderTypeOllama:
		embeddingLlm, err := ollama.New(ollama.WithModel(embeddingModel), ollama.WithServerURL(url))
		if err != nil {
			log.Fatal(err)
		}

		chatLlm, err := ollama.New(ollama.WithModel(chatModel), ollama.WithServerURL(url))
		if err != nil {
			log.Fatal(err)
		}
		return &LanchChain{embeddingModel: embeddingModel, chatLlm: chatLlm, chatModel: chatModel, scoreThreshold: scoreThreshold, temperature: temperature, embeddingLlm: embeddingLlm, kbCfg: kbCfg, log: logger, crawler: crawler.New(logger), http: util.NewHttpAccessor()}
	case data.LLMProviderTypeOpenAI:
		embeddingLlm, err := openai.New(openai.WithModel(embeddingModel), openai.WithBaseURL(url), openai.WithToken("empty"))
		if err != nil {
			log.Fatal(err)
		}

		chatLlm, err := openai.New(openai.WithModel(chatModel), openai.WithBaseURL(url), openai.WithToken("empty"))
		if err != nil {
			log.Fatal(err)
		}
		return &LanchChain{embeddingModel: embeddingModel, openAIChatLLM: chatLlm, chatModel: chatModel, scoreThreshold: scoreThreshold, temperature: temperature, openAIEmbeddingLLM: embeddingLlm, kbCfg: kbCfg, log: logger, crawler: crawler.New(logger), http: util.NewHttpAccessor()}
	}
	log.Fatal("Invalid LLM provider type")
	return nil
}
