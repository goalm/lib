package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
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

func FindMp(fileName string) (fileLoc string) {
	for _, path := range MpLocs {
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("Error reading directory " + path)
			continue
		}
		for _, file := range files {
			if file.Name() == fileName {
				fileLoc = path + "/" + fileName
				log.Println(fileName+" found at: ", fileLoc)
				return
			}
		}
	}
	fmt.Println("Model Point " + fileName + " not found in any of the paths")
	return
}

func FindFile(tbl string) (fileLoc string) {
	fileName := Conf.GetString("fileNames." + tbl)
	for _, path := range TableLocs {
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("Error reading directory " + path)
			continue
		}
		for _, file := range files {
			if file.Name() == fileName {
				fileLoc = path + "/" + fileName
				log.Println(fileName+" found at: ", fileLoc)
				return
			}
		}
	}
	fmt.Println("Table " + fileName + " not found in any of the paths")
	return
}

func GetDataFile(tbl string) string {
	return Conf.GetString("data." + tbl)
}

func GetTable(tbl string) string {
	return Conf.GetString("tables." + tbl)
}
func GetTableName(tbl string) string {
	return Conf.GetString("tableNames." + tbl)
}

func GetEnumName(enm string) string {
	return Conf.GetString("enumNames." + enm)
}

func GetFileName(s string) string {
	return Conf.GetString("Tables." + s)
}

// todo: remove this function
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
