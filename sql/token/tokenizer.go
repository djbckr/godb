package token

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

const (
	stringTypeNone = iota
	stringTypeSingle
	stringTypeDouble
	stringTypeSpecial
	stringTypeBackTk
	stringTypeBrcket
)

const (
	parenL = '('
	parenR = ')'
	braceL = '{'
	braceR = '}'
	brackL = '['
	brackR = ']'

	plus        = '+'
	minus       = '-'
	splat       = '*'
	slash       = '/'
	equals      = '='
	lessthan    = '<'
	greaterthan = '>'

	amp         = '&'
	caret       = '^'
	percent     = '%'
	pound       = '#'
	atSign      = '@'
	bang        = '!'
	tilde       = '~'
	pipe        = '|'
	backslash   = '\\'
	colon       = ':'
	semicolon   = ';'
	comma       = ','
	period      = '.'
	question    = '?'
	newLine     = '\n'
	cr          = '\r'
	tab         = '\t'
	space       = ' '
	dollar      = '$'
	quotesingle = '\''
	quotedouble = '"'
	backtick    = '`'
)

type parserData struct {
	commentDepth    int // 0 == not inside comment; positive values is depth; -1 is single-line comment
	specialQuote    *specialStringResult
	quoteType       int
	stringStartChar int
	qStrEnd         rune
	chars           []rune
	skipTo          int
	number *big.Float
}

const (
	TypeToken = iota
	TypeString
	TypePunctuation
	TypeNumber
)

type Token struct {
	Value     interface{}
	TokenType int
}

type Tokens = []Token

func Tokenize(sql string) (parsed Tokens, err error) {

	parser := &parserData{
		chars:           []rune(sql),
		stringStartChar: -1,
		quoteType:       stringTypeNone,
	}

	maxIdx := len(parser.chars) - 1

	for idx, char := range parser.chars {

		if parser.skipTo > idx {
			continue
		}

		if parser.skipTo > 0 {

			if parser.specialQuote != nil {
				parsed = append(parsed, Token{
					Value:     parser.specialQuote.contents,
					TokenType: TypeString,
				})
				parser.specialQuote = nil
				parser.quoteType = stringTypeNone
				parser.stringStartChar = -1
				parser.skipTo = 0
				continue
			} else if parser.number != nil {
				parsed = append(parsed, Token{
					Value:     parser.number,
					TokenType: TypeNumber,
				})
				parser.number = nil
				parser.stringStartChar = -1
				parser.skipTo = 0
			} else {
				parsed = append(parsed, Token{
					Value:     strings.ToUpper(string(parser.chars[parser.stringStartChar:idx])),
					TokenType: TypeToken,
				})
				parser.stringStartChar = -1
				parser.skipTo = 0
			}

		}

		var nextChar rune = 0
		var prevChar rune = 0

		if idx < maxIdx {
			nextChar = parser.chars[idx+1]
		}

		if idx > 0 {
			prevChar = parser.chars[idx-1]
		}

		switch {

		// start block level comment
		case parser.quoteType == stringTypeNone && char == slash && nextChar == splat:
			parser.commentDepth++

		// end block level comment
		case char == splat && nextChar == slash:
		// do nothing - iterate to next rune to trigger [next case]

		case parser.quoteType == stringTypeNone && prevChar == splat && char == slash:
			parser.commentDepth--

		// single line comment (allow // )
		case parser.quoteType == stringTypeNone && char == slash && nextChar == slash && parser.commentDepth == 0:
			parser.commentDepth = -1

		// single line comment ( SQL -- )
		case parser.quoteType == stringTypeNone && char == minus && nextChar == minus && parser.commentDepth == 0:
			parser.commentDepth = -1

		// end single line comment
		case parser.commentDepth == -1 && char == newLine || char == cr:
			parser.commentDepth = 0

		case parser.commentDepth != 0:
			// we'reDollarTag in the middle of a comment; skip processing now

		// start quoted identifier
		case char == quotedouble && parser.quoteType == stringTypeNone:
			parser.stringStartChar = idx + 1
			parser.quoteType = stringTypeDouble

		// start bracket quoted identifier
		case char == brackL && parser.quoteType == stringTypeNone:
			parser.stringStartChar = idx + 1
			parser.quoteType = stringTypeBrcket

		// dollar (tagged) string
		case parser.quoteType == stringTypeNone && char == dollar && checkDollarTag(parser, idx):
			parser.quoteType = stringTypeSpecial

		// start single quoted string
		case char == quotesingle && parser.quoteType == stringTypeNone:
			parser.stringStartChar = idx + 1
			parser.quoteType = stringTypeSingle

		// start back-tick string
		case char == backtick && parser.quoteType == stringTypeNone:
			parser.stringStartChar = idx + 1
			parser.quoteType = stringTypeBackTk

		// start Q-String
		case parser.quoteType == stringTypeNone && (char == 'q' || char == 'Q') && nextChar == quotesingle && checkQString(parser, idx):
			parser.quoteType = stringTypeSpecial

		// end quoted identifier
		case char == quotedouble && parser.quoteType == stringTypeDouble:
			parsed = append(parsed, Token{
				Value:     string(parser.chars[parser.stringStartChar:idx]),
				TokenType: TypeToken,
			})
			parser.stringStartChar = -1
			parser.quoteType = stringTypeNone

		// end bracket quoted identifier
		case char == brackR && parser.quoteType == stringTypeBrcket:
			parsed = append(parsed, Token{
				Value:     string(parser.chars[parser.stringStartChar:idx]),
				TokenType: TypeToken,
			})
			parser.stringStartChar = -1
			parser.quoteType = stringTypeNone

		// next char is a single-quote literal
		case char == quotesingle && parser.quoteType == stringTypeSingle && nextChar == quotesingle:
		// TODO - next is single-quote literal

		// end single quoted string
		case char == quotesingle && parser.quoteType == stringTypeSingle:
			parsed = append(parsed, Token{
				Value:     string(parser.chars[parser.stringStartChar:idx]),
				TokenType: TypeString,
			})
			parser.stringStartChar = -1
			parser.quoteType = stringTypeNone

		// end back-tick string
		case char == backtick && parser.quoteType == stringTypeBackTk:
			parsed = append(parsed, Token{
				Value:     string(parser.chars[parser.stringStartChar:idx]),
				TokenType: TypeString,
			})
			parser.stringStartChar = -1
			parser.quoteType = stringTypeNone

		case parser.quoteType != stringTypeNone:
		// we're in the middle of a string; skip processing now

		case isNumber(parser, idx):
		// isNumber will configure parser

		case allowedIdentifierChar(char):
			parser.stringStartChar = idx
			parser.skipTo = findEndOfIdentifier(parser.chars[idx:], idx)

		case
			char == parenL || char == parenR || char == plus || char == minus || char == splat ||
				char == slash || char == equals || char == lessthan || char == greaterthan ||
				char == amp || char == caret || char == percent || char == pound || char == atSign ||
				char == bang || char == tilde || char == pipe || char == backslash || char == colon ||
				char == semicolon || char == comma || char == period:
			parsed = append(parsed, Token{
				Value:     string(parser.chars[idx : idx+1]),
				TokenType: TypePunctuation,
			})

		} // end switch

	} // end loop rune range

	if parser.commentDepth > 0 {
		return nil, errors.New("unclosed comment")
	}

	if parser.specialQuote != nil {
		parsed = append(parsed, Token{
			Value:     parser.specialQuote.contents,
			TokenType: TypeString,
		})
		parser.quoteType = stringTypeNone
	} else if parser.number != nil {
		parsed = append(parsed, Token{
			Value:     parser.number,
			TokenType: TypeNumber,
		})
		parser.number = nil
	} else {
		if parser.skipTo > 0 {
			parsed = append(parsed, Token{
				Value:     strings.ToUpper(string(parser.chars[parser.stringStartChar:])),
				TokenType: TypeToken,
			})
			parser.quoteType = stringTypeNone
		}
	}

	if parser.quoteType != stringTypeNone {
		return nil, errors.New("unclosed quote")
	}

	return parsed, nil

}

var validNumRE = regexp.MustCompile(`(?i)^(?:[-+0-9]|inf)[-+._0-9a-z]*`)

func isNumber(parser *parserData, start int) bool {

	str := parser.chars[start:]
	matches := validNumRE.FindStringIndex(string(str))

	if matches == nil {
		return false
	}

	ss := string(str[matches[0]:matches[1]])
	n, _, e := big.NewFloat(0).Parse(ss, 0)

	if e != nil {
		return false
	}

	parser.skipTo = start + matches[1]
	parser.number = n

	return true

}

var reDollarTag = regexp.MustCompile(`(?mU)\$([^-$[:space:][:cntrl:]]*)\$`)
var qstrMarker = regexp.MustCompile(`(?Ui)^q'(.)`)
var qstrContent = `(?Ui)^q'%v(.*)%v'`

type specialStringResult struct {
	contents string
	end      int
}

func checkQString(parser *parserData, start int) bool {
	if parser.specialQuote = captureQString(parser.chars[start:]); parser.specialQuote != nil {
		parser.skipTo = start + parser.specialQuote.end
		return true
	}
	return false
}

func captureQString(str []rune) *specialStringResult {

	match := qstrMarker.FindAllStringSubmatch(string(str), -1)

	if match == nil {
		return nil
	}

	char1 := []rune(match[0][1])

	if len(char1) != 1 {
		panic("more than one rune??")
	}

	// clean up some regex points
	switch char1[0] {
	case '^':
		char1 = []rune("\\^")
	case '`':
		char1 = []rune("\\x60")
	case '$':
		char1 = []rune("\\$")
	case '*':
		char1 = []rune("\\*")
	case '+':
		char1 = []rune("\\+")
	case '|':
		char1 = []rune("\\|")
	case '\\':
		char1 = []rune("\\\\")
	case '?':
		char1 = []rune("\\?")
	case '.':
		char1 = []rune("\\.")
	}

	// regex doesn't like direct unicode, so encode it if necessary
	if char1[0] > 0x7e {
		char1 = []rune(fmt.Sprintf("\\x{%X}", char1[0]))
	}

	// end code will normally be start code
	char2 := char1

	// if one of [ { ( < pair with > ) } ] and fix regex problems
	switch char1[0] {
	case '[':
		char1 = []rune("\\[")
		char2 = []rune("\\]")
	case '{':
		char1 = []rune("\\{")
		char2 = []rune("\\}")
	case '(':
		char1 = []rune("\\(")
		char2 = []rune("\\)")
	case '<':
		char1 = []rune("<")
		char2 = []rune(">")
	}

	re2 := regexp.MustCompile(fmt.Sprintf(qstrContent, string(char1), string(char2)))
	cap2 := re2.FindStringSubmatchIndex(string(str))

	if cap2 == nil {
		return nil
	}

	return &specialStringResult{
		contents: string(str)[cap2[2]:cap2[3]],
		end:      cap2[1],
	}

}

func checkDollarTag(parser *parserData, start int) bool {
	if parser.specialQuote = captureDollarString(string(parser.chars)[start:]); parser.specialQuote != nil {
		parser.skipTo = start + parser.specialQuote.end
		return true
	}
	return false
}

// returns a *specialStringResult if a valid matching tag string was found; including $$$$ (empty tag, empty contents)
func captureDollarString(str string) *specialStringResult {

	rslt := reDollarTag.FindAllStringSubmatchIndex(str, -1)

	if rslt == nil {
		return nil
	}

	if len(rslt) < 2 {
		return nil
	}

	tagStr := str[rslt[0][2]:rslt[0][3]]

	i := 1

	for {

		if str[rslt[i][2]:rslt[i][3]] == tagStr {
			return &specialStringResult{
				contents: str[rslt[0][1]:rslt[i][0]],
				end:      rslt[i][3] + 1,
			}
		}

		i += 1
		if i >= len(rslt) {
			return nil
		}

	}

}

func findEndOfIdentifier(what []rune, addAmt int) (result int) {

	result = addAmt

	for _, char := range what {
		if allowedIdentifierChar(char) {
			result++
		} else {
			return
		}
	}

	return

}

func allowedIdentifierChar(char rune) bool {
	return (char >= 'A' && char <= 'Z') ||
		(char >= 'a' && char <= 'z') ||
		char == dollar || char == '_' ||
		(char >= '0' && char <= '9') ||
		char >= 192
}
