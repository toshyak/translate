package spelling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type spellingError struct {
	Pos        int      `json:"pos"`
	Len        int      `json:"len"`
	Word       string   `json:"word"`
	Suggestion []string `json:"s"`
}

// CheckSpelling checks spelling of text written in lang
// https://yandex.ru/dev/speller/doc/dg/reference/checkText-docpage/
func CheckSpelling(text string, lang string) <-chan string {
	out := make(chan string)
	go doCheck(text, lang, out)
	return out
}

func doCheck(text string, lang string, out chan string) {
	defer close(out)
	request, err := buildSpellingRequest(text, lang)
	if err != nil {
		log.Println("Cannot build request for spelling suggestions:", err)
		return
	}

	httpCallTimeout, _ := time.ParseDuration("30s")
	client := &http.Client{Timeout: httpCallTimeout}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Failed to make HTTP request to Yandex Speller:", err)
		return
	}
	defer resp.Body.Close()

	var response []spellingError
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Cannot decode response body:", err)
		return
	}

	if len(response) == 1 {
		// return all suggestions for the single error
		spellingError := response[0]
		for _, s := range spellingError.Suggestion {
			out <- strings.ReplaceAll(text, string(spellingError.Word), string(s))
		}
	} else if len(response) > 1 {
		// return only first suggestion for multiple errors
		for _, spellingError := range response {
			out <- strings.ReplaceAll(text, string(spellingError.Word), string(spellingError.Suggestion[0]))
		}
	}
}

func buildSpellingRequest(text string, lang string) (*http.Request, error) {
	requestURL := "https://speller.yandex.net/services/spellservice.json/checkText"
	// For some reasons sending POST request with JSON body like {"text": "translatethis"} doesn't work,
	// so use the same data as in https://speller.yandex.net/services/spellservice.json?op=checkText
	// Ignore capitalization errors with Options=512 https://yandex.ru/dev/speller/doc/dg/reference/speller-options-docpage/
	spellerRequest := fmt.Sprintf("text=%s&lang=%s&options=%s", url.QueryEscape(text), lang, "512")
	request, err := http.NewRequest("POST", requestURL, bytes.NewReader([]byte(spellerRequest)))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Println("Cannot create http request for url " + requestURL)
		return nil, err
	}
	return request, nil
}
