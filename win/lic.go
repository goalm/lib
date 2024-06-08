//go:build windows

package win

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	LicFile = "lic.json"
	ExpDate = "12/31/2024"
)

type User struct {
	Name string `json:"Name"`
	UUID string `json:"UUID"`
}

func ValidLicense(vaultPath string) bool {
	UUID := GetUUID()
	licFile := vaultPath + "/" + LicFile
	jsonFile, err := os.Open(licFile)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	var users []User
	err = json.NewDecoder(jsonFile).Decode(&users)
	if err != nil {
		panic(err)
	}
	stillValid := beforeExpDate()
	for _, user := range users {
		if user.UUID == UUID {
			if stillValid {
				return true
			} else {
				return false
			}
		}
	}
	fmt.Println("Please reach out to Molly / Martin to register your machine ID.")

	return false
}

func beforeExpDate() bool {
	currTime := time.Now().Unix()
	expDate, _ := time.Parse("01/02/2006", ExpDate)
	expDateStamp := expDate.Unix()

	if currTime < expDateStamp {
		log.Println(expDate.Year(), expDate.Month(), expDate.Day())
	} else {
		fmt.Println("Please reach out to Molly / Martin to request an updated version.")
		log.Println("Please reach out to Molly / Martin to request an updated version.")
		return false
	}

	return true
}
