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

// Oops, total needs to be per-prefix; needs a refactor
type suffixlist struct {
	total  int
	weight map[string]map[string]int
}

func (sl suffixlist) String() string {
	return fmt.Sprintf("%#v", sl.weight)
}

func (sl suffixlist) Inc(p, s string) {
	sl.total += 1
	premap, ok := sl.weight[p]
	if !ok {
		premap = make(map[string]int)
	}
	premap[s] += 1
	sl.weight[p] = premap
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

func makepairs(words []string) [][]string {
	ret := make([][]string, len(words)+1)
	prev := "<START>"
	for i, next := range words {
		ret[i] = []string{prev, next}
		prev = next
	}
	ret[len(ret)-1] = []string{prev, "<EOL>"}
	return ret
}

func main() {
	file, _ := os.Open(os.Args[1])
	scanner := bufio.NewScanner(file)
	stats := make(map[string]suffixlist)
	global := suffixlist{0, make(map[string]map[string]int)}
	for scanner.Scan() {
		if matched, _ := regexp.Match("^\\d\\d:\\d\\d <", scanner.Bytes()); matched {
			name, text := parseline(scanner.Text())
			text = premunge(text)
			words := strings.Split(text, " ")
			list, ok := stats[name]
			if !ok {
				list = suffixlist{0, make(map[string]map[string]int)}
			}
			pairs := makepairs(words)
			for _, pair := range pairs {
				list.Inc(pair[0], pair[1])
				global.Inc(pair[0], pair[1])
			}
			stats[name] = list
		}
	}
	fmt.Println(global)
}
