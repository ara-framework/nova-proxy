package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ara-framework/nova-proxy/logger"
)

// ViewJobError is an error happened during and after a view is requesting.
type ViewJobError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type hypernovaResult struct {
	Success bool
	Html    string
	Name    string
	Error   ViewJobError
}

type hypernovaResponse struct {
	Results map[string]hypernovaResult
}

// ModifyBody should modify specific body characters
func ModifyBody(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	logger.Fatal(err, "Cannot handle incoming html ")

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
	logger.Fatal(encodeErr, "Cannot convert batch into byte ")

	payload := string(b)

	resp, reqErr := http.Post(
		os.Getenv("HYPERNOVA_BATCH"),
		"application/json",
		strings.NewReader(payload))

	if reqErr != nil {
		logger.Error(reqErr, "Cannot reach end with html given")
		return html
	}

	defer resp.Body.Close()

	body, bodyErr := ioutil.ReadAll(resp.Body)

	if bodyErr != nil {
		log.Fatal(bodyErr)
	}

	var hypernovaResponse hypernovaResponse

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
		logger.Fatal(htmlError, "Cannot parse html element")
	}

	return html
}

func createQuery(tag string, uuid string, name string) string {
	query := fmt.Sprintf("%s[data-hypernova-id=\"%s\"][data-hypernova-key=\"%s\"]", tag, uuid, name)

	return query
}
