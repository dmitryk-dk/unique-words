package counter

import (
	"io"
	"sync"
	"unicode"

	"github.com/dmitryk-dk/unique-words/wordcount"
)

type Counter struct {
	inputStreamer wordcount.InputStream
	words         map[string]int
	wordC         chan string
	mx            sync.Mutex
}

func New(stream wordcount.InputStream) *Counter {
	c := &Counter{
		inputStreamer: stream,
		words:         make(map[string]int),
		wordC:         make(chan string),
	}
	go c.readWord()
	return c
}

func (c *Counter) readWord() {
	for word := range c.wordC {
		if word == "" {
			continue
		}
		if _, ok := c.words[word]; ok {
			c.words[word] += 1
		} else {
			c.words[word] = 1
		}
	}
	close(c.wordC)
}

func (c *Counter) WordCounts() map[string]int {
	wordCounts := make(map[string]int)
	c.mx.Lock()
	for word, count := range c.words {
		wordCounts[word] = count
	}
	c.mx.Unlock()
	return wordCounts
}

func (c *Counter) CollectWord() error {
	runes := make([]rune, 0, 100)
	for {
		r, err := c.inputStreamer.TakeChar()
		if err != nil {
			if err == io.EOF {
				c.wordC <- string(runes)
				return nil
			}
			return err
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			runes = append(runes, unicode.ToLower(r))
		} else {
			c.wordC <- string(runes)
			runes = runes[0:0]
		}
	}
}
