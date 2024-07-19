package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type KBConfig struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       struct {
		Docs []string `yaml:"docs"`
	} `yaml:"spec"`
}

func LoadKBConfig() (KBConfig, error) {
	var kbCfg KBConfig

	data, err := os.ReadFile(path.Join(configPath, kbConfigFileName))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal the data from the config file into the cfg struct
	err = yaml.Unmarshal(data, &kbCfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Loaded config: %+v\n", kbCfg)

	return kbCfg, err
}
