package synonyms

// https://yandex.ru/dev/dictionary/doc/dg/reference/lookup.html

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const apiKey = "dict.1.1.20200909T155527Z.0357057f50984cf2.ca5e23a000b42b928facbd62b402878616ad6fa3"

type synonym struct {
	Text string `json:"text"`
}

type translation struct {
	Text     string    `json:"text"`
	Synonyms []synonym `json:"syn"`
}

type definition struct {
	Translations []translation `json:"tr"`
}

type dictResponse struct {
	Definitions []definition `json:"def"`
}

// TranslateWithSynonyms returns text translation with synonyms
func TranslateWithSynonyms(text string, sourceLanguageCode string, targetLanguageCode string) ([]string, error) {
	translationDirection := fmt.Sprintf("%s-%s", sourceLanguageCode, targetLanguageCode)
	request, err := buildSpellingRequest(text, translationDirection)
	if err != nil {
		return nil, err
	}
	httpCallTimeout, _ := time.ParseDuration("30s")
	client := &http.Client{Timeout: httpCallTimeout}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Failed to make HTTP request to Yandex Dictionary:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var response dictResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Cannot decode response body:", err)
		return nil, err
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
	return translations, nil
}

func buildSpellingRequest(text string, translationDirection string) (*http.Request, error) {
	requestURL := "https://dictionary.yandex.net/api/v1/dicservice.json/lookup"
	requestParams := fmt.Sprintf("?key=%s&lang=%s&flags=12&text=%s", apiKey, translationDirection, text)
	request, err := http.NewRequest("GET", requestURL+requestParams, nil)
	if err != nil {
		log.Println("Cannot create http request for url " + requestURL)
		return nil, err
	}
	return request, nil
}
