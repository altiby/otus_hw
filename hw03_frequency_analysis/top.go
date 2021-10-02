package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func MultiReplacer(data string, symbols ...rune) string {
	builder := strings.Builder{}
	var putToBuilder bool
	for _, d := range data {
		putToBuilder = true
		for _, symbol := range symbols {
			if d == symbol {
				putToBuilder = false
				break
			}
		}
		if putToBuilder {
			builder.WriteRune(d)
		}
	}
	return builder.String()
}

// normalize преобразуем строку к нижнему регистру и удалим спецсимволы.
func normalize(data string) string {
	return MultiReplacer(strings.ToLower(data), '!', ',', '\'', '-', '.')
}

type wordStruct struct {
	Data  string
	Count int
}

func Top10(text string) []string {
	ret := make([]string, 0, 10)
	words := strings.Fields(text)
	wordsMap := make(map[string]int, len(words)/5)

	for _, word := range words {
		nword := normalize(word)
		if len(nword) == 0 {
			continue
		}
		wordsMap[nword]++
	}

	ws := make([]wordStruct, 0, len(wordsMap))
	for key, value := range wordsMap {
		ws = append(ws, wordStruct{
			Data:  key,
			Count: value,
		})
	}

	sort.Slice(ws, func(i, j int) bool {
		if ws[i].Count == ws[j].Count {
			return ws[i].Data < ws[j].Data
		}
		return ws[i].Count > ws[j].Count
	})

	for i := 0; i < 10 && i < len(ws); i++ {
		ret = append(ret, ws[i].Data)
	}

	return ret
}
