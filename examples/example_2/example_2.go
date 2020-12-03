package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dmitryk-dk/unique-words/counter"
	wordRank "github.com/dmitryk-dk/unique-words/word-rank"
	"github.com/dmitryk-dk/unique-words/wordcount"
)

func main() {
	var wg sync.WaitGroup
	var counters []*counter.Counter
	ticker := time.NewTicker(time.Second * 3)
	wordsCountC := make(chan map[string]int)
	limiter := 10

	for i := 0; i < limiter; i++ {
		counters = append(counters, counter.New(wordcount.MakeSlowReader(wordcount.ExampleText)))
	}
	for _, ctr := range counters {
		wg.Add(1)
		go func(ctr *counter.Counter) {
			if err := ctr.CollectWord(); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(ctr)
	}
	go func() {
		for range ticker.C {
			wordsCountC <- collectCounts(counters)
		}
	}()
	go func() {
		for data := range wordsCountC {
			rankings := wordRank.RankWords(data)
			for _, ranking := range rankings {
				fmt.Printf("Refreshed data: %s\n", ranking)
			}
		}
	}()
	wg.Wait()
	ticker.Stop()
}

func collectCounts(counters []*counter.Counter) map[string]int {
	for _, ctr := range counters {
		return ctr.WordCounts()
	}
	return nil
}
