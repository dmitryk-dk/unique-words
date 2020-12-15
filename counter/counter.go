package counter

import (
	"io"
	"strings"
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
		c.mx.Lock()
		if _, ok := c.words[word]; ok {
			c.words[word] += 1
		} else {
			c.words[word] = 1
		}
		c.mx.Unlock()
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

func (c *Counter) CollectWord() (err error) {
	runes := make([]rune, 0, 100)
	var r rune
	for {
		r, err = c.inputStreamer.TakeChar()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			runes = append(runes, unicode.ToLower(r))
			continue
		}
		if len(runes) > 0 {
			c.wordC <- string(runes)
			runes = runes[0:0]
		}
	}
	if len(runes) > 0 {
		c.wordC <- string(runes)
	}
	c.inputStreamer.Dispose()
	close(c.wordC)
	return
}

func (c *Counter) CollectWordWithBuilder() (err error) {
	var runes strings.Builder
	var r rune
	for {
		r, err = c.inputStreamer.TakeChar()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			runes.WriteRune(unicode.ToLower(r))
			continue
		}
		if runes.Len() > 0 {
			c.wordC <- runes.String()
			runes.Reset()
		}
	}
	if runes.Len() > 0 {
		c.wordC <- runes.String()
	}
	c.inputStreamer.Dispose()
	close(c.wordC)
	return
}

func (c *Counter) CollectWordWithoutCapacity() (err error) {
	runes := make([]rune, 0)
	var r rune
	for {
		r, err = c.inputStreamer.TakeChar()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			runes = append(runes, unicode.ToLower(r))
			continue
		}
		if len(runes) > 0 {
			c.wordC <- string(runes)
			runes = nil
		}
	}
	if len(runes) > 0 {
		c.wordC <- string(runes)
	}
	c.inputStreamer.Dispose()
	close(c.wordC)
	return
}
