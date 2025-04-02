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
	embeddingModel string
	chatModel      string
	scoreThreshold float32
	temperature    float64
	embeddingLlm   *ollama.LLM
	chatLlm        *ollama.LLM
	kbCfg          data.KBConfig
	log            *logrus.Logger
	crawler        *crawler.Crawler
	http           *util.HttpAccessor
}

func New(chatModel, embeddingModel, url string,
	scoreThreshold float32,
	temperature float64,
	kbCfg data.KBConfig,
	logger *logrus.Logger) adapters.LLMHandler {
	if len(url) == 0 {
		url = "http://127.0.0.1:11434"
	}
	embeddingLlm, err := ollama.New(ollama.WithModel(embeddingModel), ollama.WithServerURL(url))
	if err != nil {
		log.Fatal(err)
	}

	chatLlm, err := ollama.New(ollama.WithModel(chatModel), ollama.WithServerURL(url))
	if err != nil {
		log.Fatal(err)
	}

	return &LanchChain{embeddingModel: embeddingModel, chatLlm: chatLlm, chatModel: chatModel, scoreThreshold: scoreThreshold, temperature: temperature, embeddingLlm: embeddingLlm, kbCfg: kbCfg, log: logger, crawler: crawler.New(logger), http: util.NewHttpAccessor()}
}
