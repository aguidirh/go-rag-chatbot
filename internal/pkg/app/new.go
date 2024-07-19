package app

import (
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/config"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/databases/qdrant"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/llms/langchain"
	"github.com/tmc/langchaingo/embeddings"
)

type App struct {
	cfg        config.Config
	kbCfg      config.KBConfig
	LLMHandler adapters.LLMHandler
	emb        embeddings.Embedder
	VectorDB   adapters.VectorDB
}

func New(cfg config.Config, kbCfg config.KBConfig) (*App, error) {
	llmHandler := langchain.New(cfg.Spec.LLM.Model, cfg.Spec.LLM.ScoreThreshold, cfg.Spec.LLM.Temperature, kbCfg)

	emb, err := llmHandler.NewEmbedder()
	if err != nil {
		return nil, err
	}

	vectorDB := qdrant.New(cfg.Spec.VectorDB.Host, cfg.Spec.VectorDB.Port, cfg.Spec.VectorDB.Collection, emb)

	return &App{
		cfg:        cfg,
		kbCfg:      kbCfg,
		LLMHandler: llmHandler,
		emb:        emb,
		VectorDB:   vectorDB,
	}, nil
}
