package main

import (
	"log"
	"os"
	"strings"
	"sync"
	"toshyak/translate/aws"
	"toshyak/translate/spelling"
	yadictionary "toshyak/translate/yaDictionary"
	"unicode"
)

type tranaslator interface {
	Translate(string) chan string
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
	spellingCh := spelling.CheckSpelling(textToTranslate, sourceLanguage)
	awsTranslator := aws.NewTranslator(sourceLanguage, translationDirections[sourceLanguage])
	awsCh := awsTranslator.Translate(textToTranslate)
	yaDictTranslator := yadictionary.NewTranslator(sourceLanguage, translationDirections[sourceLanguage])
	yaCh := yaDictTranslator.Translate(textToTranslate)

	out := newOutput()
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		for s := range spellingCh {
			out.add(s, "", "speller", false)
		}
		wg.Done()
	}()

	go func() {
		for s := range yaCh {
			out.add(s, "", "ydict", true)
		}
		wg.Done()
	}()

	go func() {
		for s := range awsCh {
			out.add(s, "", "aws", true)
		}
		wg.Done()
	}()

	wg.Wait()
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
