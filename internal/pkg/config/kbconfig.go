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

type KbConfig struct {
	ConfigPath string
	Log        *logrus.Logger
}

func (k *KbConfig) LoadKBConfig() (data.KBConfig, error) {
	var kbCfg data.KBConfig

	data, err := os.ReadFile(path.Join(k.ConfigPath, kbConfigFileName))
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
