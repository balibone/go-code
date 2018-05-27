package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	wordsCount := make(map[string]int)
	words := strings.Fields(s)
	for _,value := range words {
		wordsCount[value]++
	}
	return wordsCount
}

func main() {
	wc.Test(WordCount)
}
