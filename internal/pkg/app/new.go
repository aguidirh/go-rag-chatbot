package app

import (
	"fmt"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/databases/qdrant"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/llms/langchain"
	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/tmc/langchaingo/embeddings"
)

type App struct {
	cfg        data.Config
	kbCfg      data.KBConfig
	LLMHandler adapters.LLMHandler
	emb        embeddings.Embedder
	VectorDB   adapters.VectorDB
}

func New(cfg data.Config, kbCfg data.KBConfig) (*App, error) {
	llmHandler := langchain.New(cfg.Spec.LLM.Model, cfg.Spec.LLM.URL, cfg.Spec.LLM.ScoreThreshold, cfg.Spec.LLM.Temperature, kbCfg)

	emb, err := llmHandler.NewEmbedder()
	if err != nil {
		return nil, err
	}

	vectorDB, err := qdrant.New(cfg.Spec.VectorDB.Host, cfg.Spec.VectorDB.Port, cfg.Spec.VectorDB.Collection, emb)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to qdrant. %v", err)
	}

	return &App{
		cfg:        cfg,
		kbCfg:      kbCfg,
		LLMHandler: llmHandler,
		emb:        emb,
		VectorDB:   vectorDB,
	}, nil
}
