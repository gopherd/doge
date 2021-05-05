// +build gengenericrbtree
//
// This program generates generic source file with custom package name,
// type of key, type of value, name of struct and comment of struct.
//
// Example usage:
//	go run -tags=gengenericrbtree gengeneric.go -pkg=orderedmap -key=int -value=string -name=OrderedMap -comment="represents an ordered map"
//
// Appends argument "-o path/to/your/filename.go" to write result to file instead of stdout.

package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

//go:embed rbtree.go
var source []byte

func main() {
	var flags struct {
		pkg     string
		key     string
		value   string
		name    string
		comment string
		output  string
	}
	flag.StringVar(&flags.pkg, "pkg", "", "Package name")
	flag.StringVar(&flags.key, "key", "", "Type of the key")
	flag.StringVar(&flags.value, "value", "", "Type of the value")
	flag.StringVar(&flags.name, "name", "", "Name of the rbtree struct")
	flag.StringVar(&flags.comment, "comment", "", "Comment of the rbtree")
	flag.StringVar(&flags.output, "o", "", "Output specified file instead of stdout")
	flag.Parse()

	const (
		prefix         = "//generic:"
		ignoredPrefix  = "ignored="
		replacePrefix  = "replace<"
		replaceSuffix  = ">"
		templatePrefix = "template<"
		templateSuffix = ">"
	)
	var (
		scanner     = bufio.NewScanner(bytes.NewBuffer(source))
		replacers   = make(map[string]string)
		lineendings = []byte{'\n'}
		buf         bytes.Buffer

		state struct {
			ignored        int
			replacing      bool
			replacePattern string
			replacement    string
		}
	)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefix) {
			// reset state
			state.ignored = 0
			state.replacing = false

			line = strings.TrimPrefix(line, prefix)
			// parse "ignored=..."
			if strings.HasPrefix(line, ignoredPrefix) {
				n, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, ignoredPrefix)))
				if err != nil {
					panic("parse ignore error: " + err.Error())
				}
				state.ignored = n
				continue
			}
			// parse "replace<...>"
			if strings.HasPrefix(line, replacePrefix) {
				s := strings.TrimSpace(strings.TrimPrefix(line, replacePrefix))
				if strings.HasSuffix(s, replaceSuffix) {
					tag := reflect.StructTag(strings.TrimSuffix(s, replaceSuffix))
					addReplacer(replacers, tag.Get("pkg"), flags.pkg)
					addReplacer(replacers, tag.Get("name"), flags.name)
					addReplacer(replacers, tag.Get("comment"), flags.comment)
					continue
				}
			}
			// parse "template<...>"
			if strings.HasPrefix(line, templatePrefix) {
				s := strings.TrimSpace(strings.TrimPrefix(line, templatePrefix))
				if strings.HasSuffix(s, templateSuffix) {
					tag := reflect.StructTag(strings.TrimSuffix(s, templateSuffix))
					if x, ok := tag.Lookup("key"); ok {
						state.replacing = true
						state.replacePattern = x
						state.replacement = flags.key
					} else if x, ok := tag.Lookup("value"); ok {
						state.replacing = true
						state.replacePattern = x
						state.replacement = flags.value
					}
					continue
				}
			}
		}
		// ignored line
		if state.ignored > 0 {
			state.ignored--
			continue
		}
		if state.replacing {
			// just replace once for each "template<...>" statement
			state.replacing = false
			line = replace(line, state.replacePattern, state.replacement)
		}
		line = replaceAll(replacers, line)
		buf.WriteString(line)
		buf.Write(lineendings)
	}
	if flags.output == "" {
		fmt.Print(buf.String())
	} else {
		if err := ioutil.WriteFile(flags.output, buf.Bytes(), 0666); err != nil {
			panic(err)
		}
	}
}

func addReplacer(replacers map[string]string, from, to string) {
	if from == "" {
		panic("addReplacer: from is empty")
	}
	if to == "" {
		panic("addReplacer: to is empty")
	}
	replacers[from] = to
}

func replaceAll(replacers map[string]string, line string) string {
	if line == "" {
		return ""
	}
	for from, to := range replacers {
		line = replace(line, from, to)
	}
	return line
}

func replace(line string, from, to string) string {
	re := regexp.MustCompile(`(^|[^_])\b` + from + `\b([^_]|$)`)
	return re.ReplaceAllString(line, `${1}`+to+`${2}`)
}
