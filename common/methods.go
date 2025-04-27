package common

import (
	"log"
	"salesmatrix/model"

	"github.com/BurntSushi/toml"
)

var (
	GlobalConfig model.Config
)

func ReadTomlFile(filename string) {

	log.Println("ReadTomlFile(+)")

	_, lErr := toml.DecodeFile(filename, &GlobalConfig)
	if lErr != nil {
		log.Fatalf("error occured in loading... (RTF01) %v", lErr)
	}

	log.Println("toml file loaded successfully !!!...")

	log.Println("ReadTomlFile(-)")

}
