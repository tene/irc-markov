package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func parseline(line string) (string, string) {
	_, rest := line[0:5], line[6:]
	vals := strings.SplitN(rest, ">", 2)
	name, text := vals[0][2:], vals[1][1:]
	return name, text
}

func breakwords(m string) {
}

type suffixlist struct {
	total  int
	weight map[string]int
}

func (sl suffixlist) String() string {
	return fmt.Sprintf("%d", sl.total)
}

func premunge(str string) string {
	str = strings.ToLower(str)
	safechars := " "
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(safechars, r) >= 0 {
			return r
		}
		if r >= 'a' && r <= 'z' {
			return r
		}
		return -1
	}, str)
}

func makepairs(words []string) []string {
	ret := make([]string, len(words)-1)
	prev := words[0]
	for i, next := range words[1:] {
		ret[i] = strings.Join([]string{prev, next}, " ")
		prev = next
	}
	return ret
}

func main() {
	file, _ := os.Open(os.Args[1])
	scanner := bufio.NewScanner(file)
	stats := make(map[string]suffixlist)
	for scanner.Scan() {
		if matched, _ := regexp.Match("^\\d\\d:\\d\\d <", scanner.Bytes()); matched {
			name, text := parseline(scanner.Text())
			text = premunge(text)
			words := strings.Split(text, " ")
			list, ok := stats[name]
			if !ok {
				list = suffixlist{0, make(map[string]int)}
			}
			pairs := makepairs(words)
			for _, pair := range pairs {
				list.total += 1
				list.weight[pair] += 1
			}
			fmt.Println(pairs)
			fmt.Println(name, list.total)
			stats[name] = list
		}
	}
	fmt.Println(stats)
}
