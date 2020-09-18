package main

import (
	"log"
	"os"
	"strings"
	"toshyak/translate/aws"
	"unicode"
)

var translationDirections = map[string]string{
	"ru": "en",
	"en": "ru",
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Pass translated text as an argument")
	}
	textToTranslate := strings.Join(os.Args[1:], " ")
	sourceLanguage := getSourceLanguage(textToTranslate)
	translatedText := aws.Translate(textToTranslate, sourceLanguage, translationDirections[sourceLanguage])
	out := newOutput()
	out.add(translatedText, "", "aws")
}

func getSourceLanguage(text string) string {
	f := func(r rune) bool {
		return unicode.Is(unicode.Cyrillic, r)
	}
	if strings.IndexFunc(text, f) != -1 {
		return "ru"
	}
	return "en"
}
