package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ara-framework/nova-proxy/logger"
	"github.com/ara-framework/nova-proxy/parser"
	env "github.com/joho/godotenv"
)

type location struct {
	Path           string
	Host           string
	ModifyResponse bool
}

type configuration struct {
	Locations []location
}

var jsonConfig configuration
var origin *url.URL

// LoadEnv should load .env file
func LoadEnv() {
	env.Load()
}

// ReadConfigFile should initialize once jsonConfig
func ReadConfigFile() {
	// logger should stop execution if there is no file found
	e, err := ioutil.ReadFile(os.Getenv("CONFIG_FILE"))
	logger.Error(err, "Config file not found")

	err = json.Unmarshal(e, &jsonConfig)
	logger.Fatal(err, "Unable to parse " + os.Getenv("CONFIG_FILE"))

}

// SetUpLocations should add handlers for config.json locations
func SetUpLocations() error {
	for _, location := range jsonConfig.Locations {
		origin, err := url.Parse(location.Host)
		logger.Fatal(err, "Malformed Host field ", location.Host)

		proxy := httputil.NewSingleHostReverseProxy(origin)
		if location.ModifyResponse {
			proxy.ModifyResponse = modifyResponse
			proxy.Director = modifyRequest(origin)
		}
		http.Handle(location.Path, proxy)
	}
	return nil
}

func isValidHeader(r *http.Response) bool {
	contentType := r.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "text/html")
}

func modifyResponse(r *http.Response) error {
	if !isValidHeader(r) {
		return nil
	}

	html, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err, "Malformed HTML in Content Body")
		return err
	}

	newHTML := parser.ModifyBody(string(html))

	r.Body = ioutil.NopCloser(strings.NewReader(newHTML))
	r.ContentLength = int64(len(newHTML))
	r.Header.Set("Content-Length", strconv.Itoa(len(newHTML)))
	return nil
}

func modifyRequest(origin *url.URL) func(req *http.Request) {
	return func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.Header.Del("Accept-Encoding")
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}
}
