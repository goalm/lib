package utils

import (
	"encoding/json"
	"os"
)

type User struct {
	Name string `json:"Name"`
	UUID string `json:"UUID"`
}

func ValidLicense(UUID string, licFile string) bool {
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

	for _, user := range users {
		if user.UUID == UUID {
			return true
		}
	}

	return false
}
