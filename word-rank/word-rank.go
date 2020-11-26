package word_rank

import (
	"fmt"
	"sort"
)

type WordRank struct {
	Word   string
	Count int
}

func (wr WordRank) String() string {
	return fmt.Sprintf("%s:%d", wr.Word, wr.Count)
}

func RankWords(words map[string]int) []WordRank {
	var counts []WordRank
	for k, v := range words {
		counts = append(counts, WordRank{k, v})
	}
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})
	return counts
}
