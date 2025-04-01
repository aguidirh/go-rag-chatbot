package app

import (
	"context"
	"fmt"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/llms/langchain"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/util"
	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/embeddings"
)

type App struct {
	cfg        data.Config
	kbCfg      data.KBConfig
	LLMHandler adapters.LLMHandler
	Embedder   embeddings.Embedder
	log        *logrus.Logger
}

func New(ctx context.Context, cfg data.Config, kbCfg data.KBConfig, skipKbLoad bool, log *logrus.Logger) (*App, error) {
	llmHandler := langchain.New(cfg.Spec.LLM.Model, cfg.Spec.LLM.EmbeddingModel, cfg.Spec.LLM.URL, cfg.Spec.LLM.ScoreThreshold, cfg.Spec.LLM.Temperature, kbCfg, log)

	embedder, err := llmHandler.NewEmbedder()
	if err != nil {
		return nil, err
	}

	if !skipKbLoad {
		kbLoader := util.NewKBLoader(ctx, &kbCfg, &cfg, llmHandler, embedder, log)
		err = kbLoader.Load()
		if err != nil {
			return nil, fmt.Errorf("failed to load KB: %v", err)
		}
	}

	return &App{
		cfg:        cfg,
		kbCfg:      kbCfg,
		LLMHandler: llmHandler,
		Embedder:   embedder,
		log:        log,
	}, nil
}
