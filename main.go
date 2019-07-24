package main

import (
	"log"
	"net/http"

	"github.com/ara-framework/nova-proxy/config"
)

func init() {
	config.LoadEnv()
	config.ReadConfigFile()
}

func main() {
	config.SetUpLocations()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
