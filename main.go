package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ViewJobError is an error happened during and after a view is requesting.
type ViewJobError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type HypernovaResult struct {
	Success bool
	Html    string
	Name    string
	Error   ViewJobError
}

type HypernovaResponse struct {
	Results map[string]HypernovaResult
}

type Location struct {
	Path           string
	Host           string
	ModifyResponse bool
}

type Configuration struct {
	Locations []Location
}

func createQuery(tag string, uuid string, name string) string {
	query := fmt.Sprintf("%s[data-hypernova-id=\"%s\"][data-hypernova-key=\"%s\"]", tag, uuid, name)

	return query
}

func modifyBody(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		log.Fatal(err)
	}

	batch := make(map[string]map[string]interface{})

	doc.Find("div[data-hypernova-key]").Each(func(i int, s *goquery.Selection) {
		uuid, uuidOk := s.Attr("data-hypernova-id")
		name, nameOk := s.Attr("data-hypernova-key")
		if !uuidOk || !nameOk {
			return
		}

		scriptQuery := createQuery("script", uuid, name)

		script := doc.Find(scriptQuery).First()

		if script == nil {
			return
		}

		content := script.Text()
		content = content[4:(len(content) - 3)]

		var data interface{}

		json.Unmarshal([]byte(content), &data)

		batch[uuid] = make(map[string]interface{})
		batch[uuid]["name"] = name
		batch[uuid]["data"] = data
	})

	if len(batch) == 0 {
		return html
	}

	b, encodeErr := json.Marshal(batch)

	if encodeErr != nil {
		log.Fatal(encodeErr)
	}

	payload := string(b)

	resp, reqErr := http.Post(os.Getenv("HYPERNOVA_BATCH"), "application/json", strings.NewReader(payload))

	if reqErr != nil {
		log.Println(reqErr)
		return html
	}

	defer resp.Body.Close()

	body, bodyErr := ioutil.ReadAll(resp.Body)

	if bodyErr != nil {
		log.Fatal(bodyErr)
	}

	var hypernovaResponse HypernovaResponse

	json.Unmarshal(body, &hypernovaResponse)

	for uuid, result := range hypernovaResponse.Results {
		divQuery := createQuery("div", uuid, result.Name)

		if !result.Success {
			doc.Find(divQuery).PrependHtml("<!-- Proxy Error: " + result.Error.Name + " -->")
			continue
		}

		scriptQuery := createQuery("script", uuid, result.Name)
		doc.Find(scriptQuery).Remove()

		doc.Find(divQuery).ReplaceWithHtml(result.Html)
	}

	html, htmlError := doc.Html()

	if htmlError != nil {
		log.Fatal(htmlError)
	}

	return html
}

func main() {
	setUpLocations()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setUpLocations() error {
	b, err := ioutil.ReadFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		fmt.Println("Config file not found")
		return err
	}

	var configuration Configuration

	json.Unmarshal(b, &configuration)

	fmt.Println(configuration)

	for _, location := range configuration.Locations {
		origin, err := url.Parse(location.Host)
		if err != nil {
			log.Fatal(err)
		} else {
			proxy := httputil.NewSingleHostReverseProxy(origin)

			if location.ModifyResponse {
				proxy.ModifyResponse = ModifyResponse
				proxy.Director = func(req *http.Request) {
					req.Header.Add("X-Forwarded-Host", req.Host)
					req.Header.Add("X-Origin-Host", origin.Host)
					req.Header.Del("Accept-Encoding")
					req.URL.Scheme = "http"
					req.URL.Host = origin.Host
				}
			}

			http.Handle(location.Path, proxy)
		}
	}

	return nil
}

func ModifyResponse(r *http.Response) error {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil
	}

	html, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
		return err
	}

	newHtml := modifyBody(string(html))

	r.Body = ioutil.NopCloser(strings.NewReader(newHtml))
	r.ContentLength = int64(len(newHtml))
	r.Header.Set("Content-Length", strconv.Itoa(len(newHtml)))
	return nil
}
