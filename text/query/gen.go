// +build ignore_genparser
//
// generate command:
//
// go run gen.go >parser.go

package main

import (
	"fmt"
	"strings"
)

const header = `// Generated by gen.go, DON NOT EDIT!
package query

import (
	"time"

	"github.com/gopherd/doge/erron"
)

type Parser struct {
	errors erron.Errors
	q      Query
}

func NewParser(q Query) *Parser {
	return &Parser{
		q: q,
	}
}

func (p *Parser) next() bool {
	return p.errors.Len() == 0
}

func (p *Parser) Err() error {
	return p.errors.All()
}`

const requiredTemplate = `
func (p *Parser) Required%s(val *%s, key string) *Parser {
	if p.next() {
		v, err := Required%s(p.q, key)
		if err != nil {
			p.errors.Append(err)
		} else {
			*val = v
		}
	}
	return p
}
`

const optionalTemplate = `
func (p *Parser) %s(val *%s, key string, dft %s) *Parser {
	if p.next() {
		v, err := %s(p.q, key, dft)
		if err != nil {
			p.errors.Append(err)
		} else {
			*val = v
		}
	}
	return p
}
`

const optionalStringTemplate = `
func (p *Parser) %s(val *%s, key string, dft %s) *Parser {
	if p.next() {
		*val = %s(p.q, key, dft)
	}
	return p
}
`

const requiredExternalTemplate = `
func (p *Parser) Required%s(val *%s.%s, key string) *Parser {
	if p.next() {
		v, err := Required%s(p.q, key)
		if err != nil {
			p.errors.Append(err)
		} else {
			*val = v
		}
	}
	return p
}
`

const optionalExternalTemplate = `
func (p *Parser) %s(val *%s.%s, key string, dft %s.%s) *Parser {
	if p.next() {
		v, err := %s(p.q, key, dft)
		if err != nil {
			p.errors.Append(err)
		} else {
			*val = v
		}
	}
	return p
}
`

func main() {
	fmt.Println(header)
	for _, t := range []string{
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"bool", "string", "float32", "float64",
	} {
		T := strings.Title(t)
		fmt.Printf(requiredTemplate, T, t, T)
		if t == "string" {
			fmt.Printf(optionalStringTemplate, T, t, t, T)
		} else {
			fmt.Printf(optionalTemplate, T, t, t, T)
		}
	}

	for _, p := range [][2]string{
		{"time", "Duration"},
	} {
		pkg := p[0]
		t := p[1]
		T := strings.Title(t)
		fmt.Printf(requiredExternalTemplate, T, pkg, t, T)
		fmt.Printf(optionalExternalTemplate, T, pkg, t, pkg, t, T)
	}
}
