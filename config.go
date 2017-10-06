package util

import (
	"encoding/json"
	"log"
	"os"
)

// ReadJSON read JSON file into a struct
func ReadJSON(jsonPath string, s interface{}) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Printf("unable to open '%s'", jsonPath)
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s)
	if err != nil {
		log.Printf("unable to decode '%s' to '%v'", jsonPath, s)
		return err
	}
	return nil
}
