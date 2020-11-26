package main

import (
	"fmt"
	"log"

	"github.com/dmitryk-dk/unique-words/counter"
	wordrank "github.com/dmitryk-dk/unique-words/word-rank"
	"github.com/dmitryk-dk/unique-words/wordcount"
)

func main() {
	content := "The cat sat on the mat."
	fastReader := wordcount.MakeFastReader(content)
	c := counter.New(fastReader)
	err := c.CollectWord()
	if err != nil {
		log.Fatal(err)
	}
	rankings := wordrank.RankWords(c.WordCounts())
	for _, ranking := range rankings {
		fmt.Println(ranking)
	}
}
