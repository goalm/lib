package utils

import (
	"log"

	"github.com/spf13/viper"
)

const (
	inputPath  = "./inputs"
	configPath = "./config"
	configType = "yaml"
)

var RunSetting *viper.Viper

func init() {
	RunSetting = viper.New()
	RunSetting.AddConfigPath(inputPath)
	RunSetting.SetConfigName(configPath)
	RunSetting.SetConfigType(configType)
	ReadConfig(RunSetting)
}

func ReadConfig(c *viper.Viper) {
	if err := c.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found")
		} else {
			log.Println("Config file found but error reading:", err)
		}
	}
}

func GetFileName(s string) string {
	return RunSetting.GetString("Tables." + s)
}
