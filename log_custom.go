package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func prepare_logger(log_filename string) {
	// Prepare Logger file
	log_file, log_err := os.OpenFile(log_filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if log_err != nil {
		log.Fatal(log_err)
	}
	defer log_file.Close()

	log.SetOutput(log_file)
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%",
	})
	log.SetLevel(log.TraceLevel)
}
