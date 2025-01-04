package logger

import (
	"log"
)

func LogError(err error) {
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
}

func LogInfo(info string) {
	log.Printf("INFO: %v\n", info)
}
