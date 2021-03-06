package counter

import (
	"testing"

	"github.com/dmitryk-dk/unique-words/wordcount"
	"github.com/stretchr/testify/assert"
)

func TestCounter_CollectWord(t *testing.T) {
	type fields struct {
		inputStreamer wordcount.InputStream
		words         map[string]int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test 1 (no new line)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader("some new words"),
				words: map[string]int{
					"some":  1,
					"new":   1,
					"words": 1,
				},
			},
			wantErr: false,
		},
		{
			name: "test 2 (new line at the end)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader("some new words\n"),
				words: map[string]int{
					"some":  1,
					"new":   1,
					"words": 1,
				},
			},
			wantErr: false,
		},
		{
			name: "test 3 (empty string)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader(""),
				words:         map[string]int{},
			},
			wantErr: false,
		},
		{
			name: "test 4 (one word)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader("one"),
				words:         map[string]int{"one": 1},
			},
			wantErr: false,
		},
		{
			name: "test 5 (empty with new line)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader(" \n"),
				words:         map[string]int{},
			},
			wantErr: false,
		},
		{
			name: "test 6 (many new lines and spaces)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader("\n one   \n"),
				words:         map[string]int{"one": 1},
			},
			wantErr: false,
		},
		{
			name: "test 7 (all letters in camel case word)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader("WORD another word"),
				words:         map[string]int{"word": 2, "another": 1},
			},
			wantErr: false,
		},
		{
			name: "test 8 (camel case word)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader("EXPECT one word"),
				words:         map[string]int{"one": 1, "expect": 1, "word": 1},
			},
			wantErr: false,
		},
		{
			name: "test 9 (camel case word)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader("EXPECT one word"),
				words:         map[string]int{"one": 1, "expect": 1, "word": 1},
			},
			wantErr: false,
		},
		{
			name: "test 10 (has punctuation)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader("EXPECT one word!!!"),
				words:         map[string]int{"one": 1, "expect": 1, "word": 1},
			},
			wantErr: false,
		},
		{
			name: "test 11 (has punctuation in different places)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader("I like, this... word EXPECT one word!!!"),
				words:         map[string]int{"expect": 1, "i": 1, "like": 1, "one": 1, "this": 1, "word": 2},
			},
			wantErr: false,
		},
		{
			name: "test 12 (just punctuation and numbers)",
			fields: fields{
				inputStreamer: wordcount.MakeFastReader(".,?!123123%&"),
				words:         map[string]int{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.fields.inputStreamer)
			if err := c.CollectWord(); (err != nil) != tt.wantErr {
				t.Errorf("CollectWord() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.fields.words, c.WordCounts())
		})
	}
}

// Test result
// BenchmarkCounter_CollectWord-8                           1439419               804 ns/op             240 B/op          4 allocs/op
func BenchmarkCounter_CollectWord(b *testing.B) {
	inputStreamer := wordcount.MakeFastReader(wordcount.ExampleText)
	for n := 0; n < b.N; n++ {
		c := New(inputStreamer)
		if err := c.CollectWord(); err != nil {
			b.Errorf("CollectWord() error = %v", err)
		}
		c.WordCounts()
	}

}

// Test result
// BenchmarkCounter_CollectWordWithBuilder-8                1499203               797 ns/op             240 B/op          4 allocs/op
func BenchmarkCounter_CollectWordWithBuilder(b *testing.B) {
	inputStreamer := wordcount.MakeFastReader(wordcount.ExampleText)
	for n := 0; n < b.N; n++ {
		c := New(inputStreamer)
		if err := c.CollectWordWithBuilder(); err != nil {
			b.Errorf("CollectWordWithBuilder() error = %v", err)
		}
		c.WordCounts()
	}
}

// Test result
// BenchmarkCounter_CollectWordWithoutCapacity-8            1493581               803 ns/op             240 B/op          4 allocs/op
func BenchmarkCounter_CollectWordWithoutCapacity(b *testing.B) {
	inputStreamer := wordcount.MakeFastReader(wordcount.ExampleText)
	for n := 0; n < b.N; n++ {
		c := New(inputStreamer)
		if err := c.CollectWordWithoutCapacity(); err != nil {
			b.Errorf("CollectWordWithBuilder() error = %v", err)
		}
		c.WordCounts()
	}
}
