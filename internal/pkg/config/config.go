package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       struct {
		VectorDB struct {
			Host       string `yaml:"host"`
			Port       string `yaml:"port"`
			Collection string `yaml:"collection"`
		} `yaml:"qdrant"`
		LLM struct {
			Model          string  `yaml:"model"`
			ScoreThreshold float32 `yaml:"scoreThreshold"`
			Temperature    float64 `yaml:"temperature"`
		} `yaml:"llm"`
		Server struct {
			Port string `yaml:"port"`
		} `yaml:"server"`
	} `yaml:"spec"`
}

func LoadConfig() (Config, error) {
	var cfg Config

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
