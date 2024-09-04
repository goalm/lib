package utils

import (
	"context"
	"github.com/spf13/viper"
	"log"
)

const (
	inputPath  = "inputs"
	configName = "config"
	configType = "yaml"
)

var ctx context.Context
var Conf *viper.Viper
var TableLocs []string
var Tables []string
var Enums []string
var EnumLocs []string
var MpLocs []string
var Mps []string

func init() {
	ctx = context.Background()
	Conf = viper.New()
	Conf.AddConfigPath(inputPath)
	Conf.AddConfigPath("../inputs") // for apps at the same level as inputs
	Conf.SetConfigName(configName)
	Conf.SetConfigType(configType)
	ReadConfig(Conf)
	TableLocs = Conf.GetStringSlice("paths.tableLocs")
	Tables = Conf.GetStringSlice("tables")
	Enums = Conf.GetStringSlice("enums")
	EnumLocs = Conf.GetStringSlice("paths.enumLocs")
	MpLocs = Conf.GetStringSlice("paths.mpLocs")
	Mps = Conf.GetStringSlice("mpFiles")
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

func GetEnumName(enm string) string {
	return Conf.GetString("enumNames." + enm)
}

func GetFileName(s string) string {
	return Conf.GetString("Tables." + s)
}
