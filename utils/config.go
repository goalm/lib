package utils

import (
	"log"

	"github.com/spf13/viper"
)

const (
	inputPath  = "inputs"
	configName = "config"
	configType = "yaml"
)

var Conf *viper.Viper

func init() {
	Conf = viper.New()
	Conf.AddConfigPath(inputPath)
	Conf.SetConfigName(configName)
	Conf.SetConfigType(configType)
	ReadConfig(Conf)
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
	return Conf.GetString("Tables." + s)
}

func GetOutPath(s string) string {
	return Conf.GetString("OutPath." + s)
}

func GetPaths(s string) string {
	return Conf.GetString("Paths." + s)
}

// For formula parser
type FormulaParserRun struct {
	Name    string
	PrdFile string
	LibFile string
}

func GetFormulaParserRuns() []FormulaParserRun {
	var runs []FormulaParserRun
	Conf.UnmarshalKey("Runs", &runs)
	return runs
}

// For Deterministic model
func GetTbl(s string) string {
	return GetPaths("inPath") + "/" + Conf.GetString("Tables."+s)
}
