package cache

import (
	"encoding/json"
	"log"
	"os"
)

func makeDump(filename string, pull any) {
	data, err := json.Marshal(pull)
	if err != nil {
		log.Println("ERROR can not marshall pull to "+filename+": ", err.Error())
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Println("ERROR can not write pull to  "+filename+": ", err.Error())
	}
	log.Println("INFO dump saved successfully to " + filename)
}
func loadFromDump(filename string, pull any) error {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &pull)
	return err
}
