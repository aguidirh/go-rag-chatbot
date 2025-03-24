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

func (c *Config) LoadConfig() (data.Config, error) {
	var cfg data.Config

	data, err := os.ReadFile(path.Join(configPath, configFileName))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal the data from the config file into the cfg struct
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Loaded config: %+v\n", cfg)

	return cfg, err
}
