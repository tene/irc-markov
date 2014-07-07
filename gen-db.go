package main

import (
	"bufio"
	"fmt"
  "os"
  "regexp"
  "strings"
)

type msg struct {
  name string
  text string
}

func (m msg) String() string {
  return fmt.Sprintf("<%s> %s", m.name, m.text)
}

func parseline(line string) msg {
  _, rest := line[0:5], line[6:]
  vals := strings.SplitN(rest, ">", 2)
  name, text := vals[0][2:], vals[1][1:]
  return msg{name, text}
}

func breakwords(m msg) {
}

func main() {
  file, _ := os.Open(os.Args[1])
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    if matched, _ := regexp.Match("^\\d\\d:\\d\\d <", scanner.Bytes()); matched {
      m := parseline(scanner.Text())
      fmt.Println(m)
    }
  }
}
