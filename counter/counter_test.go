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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.fields.inputStreamer)
			if err := c.CollectWord(); (err != nil) != tt.wantErr {
				t.Errorf("CollectWord() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.fields.words, c.words)
		})
	}
}
