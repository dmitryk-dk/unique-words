package counter

import (
	"io"
	"strings"
	"sync"
	"unicode"

	"github.com/dmitryk-dk/unique-words/wordcount"
)

const Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	newline = '\n'
	space   = ' '
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
		if _, ok := c.words[word]; ok {
			c.words[word] += 1
		} else {
			c.words[word] = 1
		}
	}
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
	runes := make([]rune, 0)
	for {
		if r, err := c.inputStreamer.TakeChar(); err == nil {
			if r == space || r == newline {
				c.wordC <- string(runes)
				runes = nil
				continue
			}
			if strings.Index(Alphabet, string(r)) > -1 {
				runes = append(runes, unicode.ToLower(r))
			}
		} else {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
