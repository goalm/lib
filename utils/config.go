package utils

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	once     sync.Once
	instance *config
)

type config struct {
	TableLoc string `yaml:"TableLoc"`
	Tables   `yaml:"Tables"`
}

type Tables struct {
	AssetsBonds    string `yaml:"AssetsBonds"`
	AssetsEquities string `yaml:"AssetsEquities"`
	AssetsCash     string `yaml:"AssetsCash"`
}

func GetConfig() *config {
	once.Do(func() {
		yamlFile, err := os.ReadFile("./inputs/config.yaml")
		if err != nil {
			log.Fatalf("Error reading YAML file: %v", err)
		}
		err = yaml.Unmarshal(yamlFile, &instance)
		if err != nil {
			log.Fatalf("Error unmarshalling YAML file: %v", err)
		}
	})
	return instance
}
