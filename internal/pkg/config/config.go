package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Log        *logrus.Logger
	ConfigPath string
}

func (c *Config) EnsureDefaults(cfg data.Config) data.Config {
	if cfg.Spec.VectorDB.Collection == "" {
		cfg.Spec.VectorDB.Collection = "default"
	}
	if cfg.Spec.VectorDB.VectorSize == 0 {
		cfg.Spec.VectorDB.VectorSize = 4096
	}
	if cfg.Spec.VectorDB.Host == "" {
		cfg.Spec.VectorDB.Host = "localhost"
	}
	if cfg.Spec.VectorDB.Port == "" {
		cfg.Spec.VectorDB.Port = "6333"
	}
	if cfg.Spec.VectorDB.Distance == "" {
		cfg.Spec.VectorDB.Distance = "Cosine"
	}
	if cfg.Spec.LLM.ChatModel == "" {
		cfg.Spec.LLM.ChatModel = "llama2"
	}
	if cfg.Spec.LLM.ScoreThreshold == 0.0 {
		cfg.Spec.LLM.ScoreThreshold = 0.5
	}
	if cfg.Spec.LLM.Temperature == 0.0 {
		cfg.Spec.LLM.Temperature = 0.8
	}
	if len(cfg.Spec.LLM.ProviderType) == 0 {
		cfg.Spec.LLM.ProviderType = data.LLMProviderTypeOllama
	}
	if cfg.Spec.LLM.URL == "" {
		switch cfg.Spec.LLM.ProviderType {
		case data.LLMProviderTypeOllama:
			cfg.Spec.LLM.URL = "http://localhost:11434"
		case data.LLMProviderTypeOpenAI:
			cfg.Spec.LLM.URL = "http://localhost:1234"
		default:
			cfg.Spec.LLM.URL = "http://localhost:11434"
		}

	}

	return cfg
}
func (c *Config) LoadConfig() (data.Config, error) {
	var cfg data.Config

	data, err := os.ReadFile(path.Join(c.ConfigPath, configFileName))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal the data from the config file into the cfg struct
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Loaded config: %+v\n", cfg)

	return c.EnsureDefaults(cfg), err
}
