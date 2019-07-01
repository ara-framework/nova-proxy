package logger

import (
	"io"
	"log"
)

// Warning initialize logger package
func Warning(err error, description ...string) {
	var warningHandle io.Writer
	if err != nil {
		for _, obj := range description {
			log.New(warningHandle,
				"WARNING: "+obj+" ",
				log.Ldate|log.Ltime|log.Lshortfile)
		}
	}
}

// Error initialize logger package
func Error(err error, description ...string) {
	var errorHandle io.Writer
	if err != nil {
		for _, obj := range description {
			log.New(
				errorHandle,
				"ERROR: "+obj+" ",
				log.Ldate|log.Ltime|log.Lshortfile)
		}
	}
}

// Fatal initialize logger package
func Fatal(err error, description ...string) {
	var errorHandle io.Writer
	if err != nil {
		log.Fatal(errorHandle,
			"FATAL: ",
			log.Ldate|log.Ltime|log.Lshortfile, description, err)
	}
}
