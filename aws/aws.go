package aws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

const awsRegion = "eu-central-1"
const awsService = "translate"

type awsTranslateRequest struct {
	SourceLanguageCode string
	TargetLanguageCode string
	Text               string
}

type awsTranslateResponse struct {
	TranslatedText string
}

// Translate returns translated text from AWS Translate service
func Translate(text string) string {
	request, err := buildAndSignHTTPRequest(text)
	if err != nil {
		log.Fatal("Cannot build HTTP request. ", err)
	}

	httpCallTimeout, _ := time.ParseDuration("30s")
	client := &http.Client{Timeout: httpCallTimeout}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("Failed to make HTTP request to AWS translate API. ", err)
	}
	defer resp.Body.Close()
	var response awsTranslateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal("Cannot decode response body", err)
	}
	return response.TranslatedText
}

func buildAndSignHTTPRequest(text string) (*http.Request, error) {
	url := fmt.Sprintf("https://%s.%s.amazonaws.com/", awsService, awsRegion)
	body := awsTranslateRequest{SourceLanguageCode: "en", TargetLanguageCode: "ru", Text: text}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Println("Cannot marshal request body: " + body.Text)
		return nil, err
	}
	bodyReader := bytes.NewReader(bodyJSON)
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Println("Cannot create http request for url " + url)
		return nil, err
	}
	err = signHTTPRequest(request, bodyReader)
	if err != nil {
		log.Println("Cannot sign HTTP request", err)
		return nil, err
	}
	return request, nil
}

func signHTTPRequest(request *http.Request, body io.ReadSeeker) error {
	creds := credentials.NewEnvCredentials()
	requestSigner := v4.NewSigner(creds)

	signatureTime := time.Now()
	signatureTimeFormatted := signatureTime.Format("20060102T150405Z")
	// add headers for translate service according to https://docs.aws.amazon.com/translate/latest/dg/API_Reference.html
	request.Header.Add("Content-Type", "application/x-amz-json-1.1")
	request.Header.Add("X-Amz-Date", signatureTimeFormatted)
	request.Header.Add("X-Amz-Target", "AWSShineFrontendService_20170701.TranslateText")
	_, err := requestSigner.Sign(request, body, awsService, awsRegion, signatureTime)
	if err != nil {
		return err
	}
	return nil
}