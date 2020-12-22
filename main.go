package main

import (
	"log"
	"os"
	"strings"
	"toshyak/translate/aws"
	"toshyak/translate/spelling"
	yadictionary "toshyak/translate/yaDictionary"
	"unicode"
)

type tranaslator interface {
	Translate(string) ([]string, error)
}

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
	spellingSuggestions, err := spelling.CheckSpelling(textToTranslate, sourceLanguage)
	if err != nil {
		log.Println("Failed to check spelling", err)
	}
	awsTranslator := aws.NewTranslator(sourceLanguage, translationDirections[sourceLanguage])
	translatedText, err := awsTranslator.Translate(textToTranslate)
	yaDictTranslator := yadictionary.NewTranslator(sourceLanguage, translationDirections[sourceLanguage])
	translatedTextWithSynonyms, err := yaDictTranslator.Translate(textToTranslate)
	out := newOutput()
	for _, s := range spellingSuggestions {
		out.add(s, "", "speller", false)
	}
	for _, s := range translatedText {
		out.add(s, "", "aws", true)
	}
	for _, s := range translatedTextWithSynonyms {
		out.add(s, "", "ydict", true)
	}
	out.print()
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
