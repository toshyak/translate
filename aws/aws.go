package aws

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	creds := credentials.NewEnvCredentials()
	requestSigner := v4.NewSigner(creds)

	url := fmt.Sprintf("https://%s.%s.amazonaws.com/", awsService, awsRegion)
	body := awsTranslateRequest{SourceLanguageCode: "en", TargetLanguageCode: "ru", Text: text}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Fatal("Cannot marshal request body: " + body.Text)
	}
	bodyReader := bytes.NewReader(bodyJSON)
	request, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		log.Fatal("Cannot create http request for url " + url)
	}
	signatureTime := time.Now()
	signatureTimeFormatted := signatureTime.Format("20060102T150405Z")
	// add headers for translate service according to https://docs.aws.amazon.com/translate/latest/dg/API_Reference.html
	request.Header.Add("Content-Type", "application/x-amz-json-1.1")
	request.Header.Add("X-Amz-Date", signatureTimeFormatted)
	request.Header.Add("X-Amz-Target", "AWSShineFrontendService_20170701.TranslateText")
	_, err = requestSigner.Sign(request, bodyReader, awsService, awsRegion, signatureTime)
	if err != nil {
		log.Fatal("Cannot sign a request")
	}
	client := &http.Client{}
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
