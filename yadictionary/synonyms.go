package yadictionary

// https://yandex.ru/dev/dictionary/doc/dg/reference/lookup.html

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const apiKey = "dict.1.1.20200909T155527Z.0357057f50984cf2.ca5e23a000b42b928facbd62b402878616ad6fa3"

type translation struct {
	Text     string `json:"text"`
	Synonyms []struct {
		Text string `json:"text"`
	} `json:"syn"`
}

type dictResponse struct {
	Definitions []struct {
		Translations []translation `json:"tr"`
	} `json:"def"`
}

// Translator object
type Translator struct {
	sourceLanguageCode string
	targetLanguageCode string
}

// NewTranslator creates new translator object
func NewTranslator(sourceLanguageCode string, targetLanguageCode string) *Translator {
	t := Translator{sourceLanguageCode, targetLanguageCode}
	return &t
}

// Translate returns text translation with synonyms
func (t *Translator) Translate(text string) <-chan string {
	out := make(chan string)
	go doTranslate(*t, text, out)
	return out
}

func doTranslate(t Translator, text string, out chan string) {
	defer close(out)
	translationDirection := fmt.Sprintf("%s-%s", t.sourceLanguageCode, t.targetLanguageCode)
	request, err := buildSpellingRequest(text, translationDirection)
	if err != nil {
		return
	}
	httpCallTimeout, _ := time.ParseDuration("30s")
	client := &http.Client{Timeout: httpCallTimeout}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Failed to make HTTP request to Yandex Dictionary:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Got a non 200 HTTP response code from Yandex Dictionary")
		log.Println("Response code received:", resp.StatusCode)
		return
	}

	var response dictResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Cannot decode response body:", err)
		return
	}

	var translations []string

	for _, def := range response.Definitions {
		for _, tr := range def.Translations {
			translations = append(translations, tr.Text)
			for _, syn := range tr.Synonyms {
				translations = append(translations, syn.Text)
			}
		}
	}
	for _, s := range translations {
		out <- s
	}
}

func buildSpellingRequest(text string, translationDirection string) (*http.Request, error) {
	requestURL := "https://dictionary.yandex.net/api/v1/dicservice.json/lookup"
	v := url.Values{}
	v.Add("key", apiKey)
	v.Add("lang", translationDirection)
	v.Add("flags", "12")
	v.Add("text", text)
	request, err := http.NewRequest("GET", requestURL+"?"+v.Encode(), nil)
	if err != nil {
		log.Println("Cannot create http request for url " + requestURL)
		return nil, err
	}
	return request, nil
}
