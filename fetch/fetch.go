package fetch

import (
	"encoding/json"
	"net/http"
	"strings"
)

const WIKIPEDIA_API_URL = "https://en.wikipedia.org/w/api.php"

func ExtractPageName(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func FetchLinksFromWikipedia(pageName string) ([]string, error) {
	params := map[string]string{
		"action":  "query",
		"format":  "json",
		"prop":    "links",
		"titles":  pageName,
		"pllimit": "max",
	}

	var result map[string]interface{}
	if err := httpGet(params, &result); err != nil {
		return nil, err
	}

	pages, ok := result["query"].(map[string]interface{})["pages"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	var links []string
	for _, page := range pages {
		pageData, ok := page.(map[string]interface{})
		if !ok {
			continue
		}

		for _, link := range pageData["links"].([]interface{}) {
			linkData, ok := link.(map[string]interface{})
			if !ok {
				continue
			}

			if title, ok := linkData["title"].(string); ok {
				if ns, ok := linkData["ns"].(float64); ok && ns == 0 {
					links = append(links, title)
				}
			}
		}
	}

	return links, nil
}

func httpGet(params map[string]string, result interface{}) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", WIKIPEDIA_API_URL, nil)

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(result)
}
