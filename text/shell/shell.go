/*
Copyright 2012 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

NOTE: Fixed by gopherd.
*/

package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// TokenType is a top-level token classification: A word, space, comment, unknown.
type TokenType int

// runeTokenClass is the type of a UTF-8 character classification: A quote, space, escape.
type runeTokenClass int

// the internal state used by the lexer state machine
type lexerState int

// Named classes of UTF-8 runes
const (
	spaceRunes            = " \t\r\n"
	escapingQuoteRunes    = `"`
	nonEscapingQuoteRunes = "'"
	escapeRunes           = `\`
	commentRunes          = "#"
)

// Classes of rune token
const (
	unknownRuneClass runeTokenClass = iota
	spaceRuneClass
	escapingQuoteRuneClass
	nonEscapingQuoteRuneClass
	escapeRuneClass
	commentRuneClass
	eofRuneClass
)

// Classes of lexographic token
const (
	UnknownToken TokenType = iota
	WordToken
	SpaceToken
	CommentToken
)

// Lexer state machine states
const (
	startState           lexerState = iota // no runes have been seen
	inWordState                            // processing regular runes in a word
	escapingState                          // we have just consumed an escape rune; the next rune is literal
	escapingQuotedState                    // we have just consumed an escape rune within a quoted string
	quotingEscapingState                   // we are within a quoted string that supports escaping ("...")
	quotingState                           // we are within a string that does not support escaping ('...')
	commentState                           // we are within a comment (everything following an unquoted or unescaped #
)

// tokenClassifier is used for classifying rune characters.
type tokenClassifier map[rune]runeTokenClass

func (typeMap tokenClassifier) addRuneClass(runes string, tokenType runeTokenClass) {
	for _, runeChar := range runes {
		typeMap[runeChar] = tokenType
	}
}

// newDefaultClassifier creates a new classifier for ASCII characters.
func newDefaultClassifier() tokenClassifier {
	t := tokenClassifier{}
	t.addRuneClass(spaceRunes, spaceRuneClass)
	t.addRuneClass(escapingQuoteRunes, escapingQuoteRuneClass)
	t.addRuneClass(nonEscapingQuoteRunes, nonEscapingQuoteRuneClass)
	t.addRuneClass(escapeRunes, escapeRuneClass)
	t.addRuneClass(commentRunes, commentRuneClass)
	return t
}

// ClassifyRune classifiees a rune
func (t tokenClassifier) ClassifyRune(runeVal rune) runeTokenClass {
	return t[runeVal]
}

// Lexer turns an input stream into a sequence of tokens. Whitespace and comments are skipped.
type Lexer Tokenizer

// NewLexer creates a new lexer from an input stream.
func NewLexer(r io.Reader) *Lexer {
	return (*Lexer)(NewTokenizer(r))
}

// Next returns the next word, or an error. If there are no more words,
// the error will be io.EOF.
func (l *Lexer) Next() ([]rune, bool, error) {
	for {
		tokenType, value, end, err := (*Tokenizer)(l).Next()
		if err != nil {
			return nil, end, err
		}
		switch tokenType {
		case WordToken:
			return value, end, nil
		case CommentToken:
			return nil, end, nil
		default:
			return nil, end, fmt.Errorf("Unknown token type: %v", tokenType)
		}
	}
}

// Tokenizer turns an input stream into a sequence of typed tokens
type Tokenizer struct {
	input      io.RuneReader
	classifier tokenClassifier
}

// NewTokenizer creates a new tokenizer from an input stream.
func NewTokenizer(r io.Reader) *Tokenizer {
	t := &Tokenizer{
		classifier: newDefaultClassifier(),
	}
	if rr, ok := r.(io.RuneReader); ok {
		t.input = rr
	} else {
		t.input = bufio.NewReader(r)
	}
	return t
}

// scanStream scans the stream for the next token using the internal state machine.
// It will panic if it encounters a rune which it does not know how to handle.
func (t *Tokenizer) scanStream() (tokenType TokenType, value []rune, end bool, err error) {
	state := startState
	var nextRune rune
	var nextRuneType runeTokenClass

	for {
		nextRune, _, err = t.input.ReadRune()
		nextRuneType = t.classifier.ClassifyRune(nextRune)

		if err == io.EOF {
			nextRuneType = eofRuneClass
			err = nil
		} else if err != nil {
			return
		}

		if nextRune == '\r' {
			nextRune, _, err = t.input.ReadRune()
			nextRuneType = t.classifier.ClassifyRune(nextRune)
			if err == io.EOF {
				nextRuneType = eofRuneClass
				err = nil
			} else if err != nil {
				return
			}
			if nextRune != '\n' {
				err = errors.New("newline expected after escape CR")
				return
			}
		}

		switch state {
		case startState: // no runes read yet
			{
				switch nextRuneType {
				case eofRuneClass:
					err = io.EOF
					return
				case spaceRuneClass:
				case escapingQuoteRuneClass:
					tokenType = WordToken
					state = quotingEscapingState
				case nonEscapingQuoteRuneClass:
					tokenType = WordToken
					state = quotingState
				case escapeRuneClass:
					tokenType = WordToken
					state = escapingState
				case commentRuneClass:
					tokenType = CommentToken
					state = commentState
				default:
					tokenType = WordToken
					value = append(value, nextRune)
					state = inWordState
				}
			}
		case inWordState: // in a regular word
			{
				switch nextRuneType {
				case eofRuneClass:
					return
				case spaceRuneClass:
					if nextRune == '\n' {
						if value == nil {
							value = make([]rune, 0)
						}
						end = true
					}
					return
				case escapingQuoteRuneClass:
					state = quotingEscapingState
				case nonEscapingQuoteRuneClass:
					state = quotingState
				case escapeRuneClass:
					state = escapingState
				default:
					value = append(value, nextRune)
				}
			}
		case escapingState: // the rune after an escape character
			{
				switch nextRuneType {
				case eofRuneClass:
					err = errors.New("EOF found after escape character")
					return
				default:
					if nextRune == '\n' {
						return
					}
					state = inWordState
					value = append(value, nextRune)
				}
			}
		case escapingQuotedState: // the next rune after an escape character, in double quotes
			{
				switch nextRuneType {
				case eofRuneClass:
					err = errors.New("EOF found after escape character")
					return
				default:
					state = quotingEscapingState
					value = append(value, nextRune)
				}
			}
		case quotingEscapingState: // in escaping double quotes
			{
				switch nextRuneType {
				case eofRuneClass:
					err = errors.New("EOF found when expecting closing quote")
					return
				case escapingQuoteRuneClass:
					state = inWordState
				case escapeRuneClass:
					state = escapingQuotedState
				default:
					value = append(value, nextRune)
				}
			}
		case quotingState: // in non-escaping single quotes
			{
				switch nextRuneType {
				case eofRuneClass:
					err = errors.New("EOF found when expecting closing quote")
					return
				case nonEscapingQuoteRuneClass:
					state = inWordState
				default:
					value = append(value, nextRune)
				}
			}
		case commentState: // in a comment
			{
				switch nextRuneType {
				case eofRuneClass:
					return
				case spaceRuneClass:
					if nextRune == '\n' {
						state = startState
						return
					} else {
						value = append(value, nextRune)
					}
				default:
					value = append(value, nextRune)
				}
			}
		default:
			err = fmt.Errorf("Unexpected state: %v", state)
			return
		}
	}
}

// Next returns the next token in the stream.
func (t *Tokenizer) Next() (TokenType, []rune, bool, error) {
	return t.scanStream()
}

// Split partitions a string into a slice of strings.
func Split(s string) ([]string, error) {
	l := NewLexer(strings.NewReader(s))
	subStrings := make([]string, 0)
	for {
		word, _, err := l.Next()
		if err != nil {
			if err == io.EOF {
				return subStrings, nil
			}
			return subStrings, err
		}
		if word != nil {
			subStrings = append(subStrings, string(word))
		}
	}
}
