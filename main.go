package main

import (
	"log"
	"os"
	"strings"
	"toshyak/translate/aws"
	"toshyak/translate/spelling"
	"toshyak/translate/synonyms"
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
	spellingSuggestions, err := spelling.CheckSpelling(textToTranslate, sourceLanguage)
	if err != nil {
		log.Println("Failed to check spelling", err)
	}
	translatedText := aws.Translate(textToTranslate, sourceLanguage, translationDirections[sourceLanguage])
	translatedTextWithSynonyms, err := synonyms.TranslateWithSynonyms(textToTranslate, sourceLanguage, translationDirections[sourceLanguage])
	out := newOutput()
	for _, s := range spellingSuggestions {
		out.add(s, "", "speller", false)
	}
	out.add(translatedText, "", "aws", true)
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
