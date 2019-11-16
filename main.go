package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ara-framework/nova-proxy/config"
	"github.com/gookit/color"
)

func init() {
	config.LoadEnv()
	config.ReadConfigFile()
}

func main() {
	config.SetUpLocations()

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "8080"
	}

	color.Info.Printf("Nova proxy running on http://0.0.0.0:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
