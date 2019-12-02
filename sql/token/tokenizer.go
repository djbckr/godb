package token

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strings"
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

const (
	TypeToken = iota
	TypeString
	TypePunctuation
	TypeNumber
	TypeComment
	TypeHint
)

type Token struct {
	Value     interface{}
	TokenType int
}

type Tokens = []*Token

type tokenizer struct {
	idx    int    // current string pointer
	maxIdx int    // total length of string - 1
	chars  string // the unchanged string
	slice  string // points to string[idx:]
	tokens Tokens // result tokens
}

func Tokenize(sql string) (parsed Tokens, err error) {

	tdata := tokenizer{
		idx:    0,
		chars:  sql,
		maxIdx: len(sql) - 1,
	}

	for tdata.idx < tdata.maxIdx {

		tdata.slice = tdata.chars[tdata.idx:]
		switch {
		case processRE(&tdata, lineHintRE, TypeHint, false):
		case processRE(&tdata, blockHintRE, TypeHint, false):
		case processRE(&tdata, lineCommentRE, TypeComment, false):
		case processRE(&tdata, blockCommentRE, TypeComment, false):
		case captureString(&tdata, captureQString):
		case captureString(&tdata, captureDollarString):
		case captureString(&tdata, captureSingleQuote):
		case processRE(&tdata, tokenRE, TypeToken, true):
		case processRE(&tdata, doubleQuoteRE, TypeToken, false):
		case processRE(&tdata, brackQuoteRE, TypeToken, false):
		case processNumber(&tdata):
		case processPunctuation(&tdata):
		default:
			if len(tdata.slice) > 30 {
				tdata.slice = tdata.slice[:30]
			}
			return nil, errors.New(fmt.Sprintf("Invalid SQL text found near: %v", tdata.slice))
		}
	}

	return tdata.tokens, nil
}

// In all below RE's, skip leading space; makes it easy to skip white space where needed

// find /*+ capture this */ ; . includes newline
var blockHintRE = regexp.MustCompile(`(?s)^[[:space:]]*?/\*\+(.*)\*/`)
// find --+ capture this to end of line
var lineHintRE = regexp.MustCompile(`(?s)^[[:space:]]*--\+([ \t\S]*)(?:[\n\r])?`)
// find /* capture this */ ; . includes newline
var blockCommentRE = regexp.MustCompile(`(?s)^[[:space:]]*?/\*(.*)\*/`)
// find -- capture this to end of line
var lineCommentRE = regexp.MustCompile(`(?s)^[[:space:]]*--([ \t\S]*)(?:[\n\r])?`)
// Unquoted (double, bracket) tokens are case-insensitive: first letter may be $, A-Z, and any unicode character
// following letters include the above, plus _, and 0-9
var tokenRE = regexp.MustCompile(`(?i)^[[:space:]]*([$A-Z\x{0080}-\x{FFEE}][$_A-Z0-9\x{0080}-\x{FFEE}]*)`)
// Capture double-quote token.  Cannot include a double-quote inside, and newlines are not allowed
var doubleQuoteRE = regexp.MustCompile(`(?U)^[[:space:]]*"([ !\x{0023}-\x{007E}\x{0080}-\x{FFEE}]*)"`)
// Capture bracket-quote token. Cannot include [ and ] and newlines are not allowed
var brackQuoteRE = regexp.MustCompile(`(?U)^[[:space:]]*\[([\x{0020}-\x{005A}\x{005C}\x{005E}-\x{007E}\x{0080}-\x{FFEE}]*)]`)
// Capture any single graph character
var punctuationRE = regexp.MustCompile(`(?i)^[[:space:]]*([[:graph:]])`)

func processRE(tdata *tokenizer, re *regexp.Regexp, typ int, upper bool) bool {

	rslt := re.FindStringSubmatchIndex(tdata.slice)

	if rslt != nil {
		str := tdata.slice[rslt[2]:rslt[3]]
		if upper {
			str = strings.ToUpper(str)
		}
		tdata.tokens = append(tdata.tokens, &Token{
			Value:     str,
			TokenType: typ,
		})
		tdata.idx += rslt[1]
		return true
	}

	return false
}

type specStrFn = func(string) *stringCaptureResult

func captureString(tdata *tokenizer, fn specStrFn) bool {
	rslt := fn(tdata.slice)
	if rslt != nil {
		tdata.tokens = append(tdata.tokens, &Token{
			Value:     rslt.contents,
			TokenType: TypeString,
		})
		tdata.idx += rslt.end
		return true
	}
	return false
}

func processPunctuation(tdata *tokenizer) bool {

	matches := punctuationRE.FindStringSubmatchIndex(tdata.slice)

	if matches == nil {
		return false
	}

	char := tdata.slice[matches[2]]

	if char == parenL || char == parenR || char == plus || char == minus || char == splat ||
		char == slash || char == equals || char == lessthan || char == greaterthan ||
		char == amp || char == caret || char == percent || char == pound || char == atSign ||
		char == bang || char == tilde || char == pipe || char == backslash || char == colon ||
		char == semicolon || char == comma || char == period {
		tdata.tokens = append(tdata.tokens, &Token{
			Value:     string(char),
			TokenType: TypePunctuation,
		})
		tdata.idx += matches[3]
		return true
	}

	return false
}

var validNumRE = regexp.MustCompile(`(?i)^[[:space:]]*((?:[-+0-9]|inf)[-+._0-9a-z]*)`)

func processNumber(tdata *tokenizer) bool {

	matches := validNumRE.FindStringSubmatchIndex(tdata.slice)

	if matches == nil {
		return false
	}

	ss := tdata.slice[matches[2]:matches[3]]
	n, _, e := big.NewFloat(0).Parse(ss, 0)

	if e != nil {
		return false
	}

	tdata.tokens = append(tdata.tokens, &Token{
		Value:     n,
		TokenType: TypeNumber,
	})
	tdata.idx += matches[3]

	return true
}

type stringCaptureResult struct {
	contents string
	end      int
}

func captureSingleQuote(str string) *stringCaptureResult {

	i := 0

	reader := strings.NewReader(str)

	var result []rune

	// position to just after quote char
	for {

		c, num, err := reader.ReadRune()

		if err == io.EOF {
			return nil
		}

		i += num

		if c == quotesingle {
			break
		}

		if c == 32 || c == 13 || c == 10 || c == 9 {
			continue
		}

		return nil
	}

	for {

		c, num, err := reader.ReadRune()

		if err == io.EOF {
			return nil
		}

		i += num

		if c == quotesingle {
			return &stringCaptureResult{
				contents: string(result),
				end:      i,
			}
		}

		if c == backslash {
			// process escape char; move to next char
			c, num, err = reader.ReadRune()

			if err == io.EOF {
				return nil
			}

			i += num

			switch c {
			case backslash:
				result = append(result, backslash)
			case 'n':
				result = append(result, 10)
			case 'r':
				result = append(result, 13)
			case 't':
				result = append(result, 9)
			case quotesingle:
				result = append(result, quotesingle)
			}
		} else {
			result = append(result, c)
		}

	}

}

var qstrMarker = regexp.MustCompile(`(?Ui)^[[:space:]]*q'(.)`)
var qstrContent = `(?isU)^[[:space:]]*q'%v(.*)%v'`

func captureQString(str string) *stringCaptureResult {

	match := qstrMarker.FindAllStringSubmatch(str, -1)

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

	return &stringCaptureResult{
		contents: str[cap2[2]:cap2[3]],
		end:      cap2[1],
	}

}

var reDollarTag = regexp.MustCompile(`(?mU)\$([^-$[:space:][:cntrl:]]*)\$`)

func captureDollarString(str string) *stringCaptureResult {

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
			return &stringCaptureResult{
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
