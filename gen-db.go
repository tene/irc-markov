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

type markov struct {
  name string
  suffix map[string]suffixlist
}

type suffixlist struct {
	total  int
	weight map[string]int
}

func (sl suffixlist) String() string {
  ret := ""
  for suffix, weight := range sl.weight {
    ret += fmt.Sprintf("  %s: %d\n", suffix, weight)
  }
  return ret;
}

func (m markov) String() string {
  ret := ""
  for prefix, suffixes := range m.suffix {
    ret += fmt.Sprintf("%s:\n%s\n", prefix, suffixes.String())
  }
  return ret;
}

func NewSuffixlist() suffixlist {
  return suffixlist{weight: make(map[string]int)}
}

func NewMarkov(name string) markov {
  return markov{name: name, suffix: make(map[string]suffixlist)}
}

func (m markov) Inc(prefix, suffix string) {
	sl, ok := m.suffix[prefix]
	if !ok {
		sl = NewSuffixlist()
	}
  sl.Inc(suffix)
  m.suffix[prefix] = sl
}

func (sl suffixlist) Inc(suffix string) {
	sl.total += 1
  sl.weight[suffix] += 1
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
	stats := make(map[string]markov)
	global := NewMarkov("global")
	for scanner.Scan() {
		if matched, _ := regexp.Match("^\\d\\d:\\d\\d <", scanner.Bytes()); matched {
			name, text := parseline(scanner.Text())
			text = premunge(text)
			words := strings.Split(text, " ")
			list, ok := stats[name]
			if !ok {
				list = NewMarkov(name)
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
