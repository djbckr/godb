package token

import (
	"errors"
	"regexp"
	"strings"
)

const (
	stringTypeNone = iota
	stringTypeSingle
	stringTypeDouble
	stringTypeDollar
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
	dollarQuote     *tagResult
	quoteType       int
	stringStartChar int
	chars           []rune
	skipTo          int
}

func Tokenize(sql string) (parsed []string, err error) {

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

			if parser.dollarQuote != nil {
				parsed = append(parsed, parser.dollarQuote.contents)
				parser.dollarQuote = nil
				parser.quoteType = stringTypeNone
				parser.stringStartChar = -1
				parser.skipTo = 0
				continue
			} else {
				parsed = append(parsed, strings.ToUpper(sql[parser.stringStartChar:idx]))
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
			// we're in the middle of a comment; skip processing now

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
			parser.quoteType = stringTypeDollar

		// start single quoted string
		case char == quotesingle && parser.quoteType == stringTypeNone:
			parser.stringStartChar = idx
			parser.quoteType = stringTypeSingle

		// start back-tick string
		case char == backtick && parser.quoteType == stringTypeNone:
			parser.stringStartChar = idx + 1
			parser.quoteType = stringTypeBackTk

		// end quoted identifier
		case char == quotedouble && parser.quoteType == stringTypeDouble:
			parsed = append(parsed, sql[parser.stringStartChar:idx])
			parser.stringStartChar = -1
			parser.quoteType = stringTypeNone

		// end bracket quoted identifier
		case char == brackR && parser.quoteType == stringTypeBrcket:
			parsed = append(parsed, sql[parser.stringStartChar:idx])
			parser.stringStartChar = -1
			parser.quoteType = stringTypeNone

		// next char is a single-quote literal
		case char == quotesingle && parser.quoteType == stringTypeSingle && nextChar == quotesingle:
		// TODO - next is single-quote literal

		// end single quoted string
		case char == quotesingle && parser.quoteType == stringTypeSingle:
			parsed = append(parsed, sql[parser.stringStartChar:idx+1])
			parser.stringStartChar = -1
			parser.quoteType = stringTypeNone

		// end back-tick string
		case char == backtick && parser.quoteType == stringTypeBackTk:
			parsed = append(parsed, sql[parser.stringStartChar:idx])
			parser.stringStartChar = -1
			parser.quoteType = stringTypeNone

		case parser.quoteType != stringTypeNone:
		// we're in the middle of a string; skip processing now

		case allowedIdentifierChar(char):
			parser.stringStartChar = idx
			parser.skipTo = findEndOfIdentifier(parser.chars[idx:], idx)

		case
			char == parenL || char == parenR || char == plus || char == minus || char == splat ||
				char == slash || char == equals || char == lessthan || char == greaterthan ||
				char == amp || char == caret || char == percent || char == pound || char == atSign ||
				char == bang || char == tilde || char == pipe || char == backslash || char == colon ||
				char == semicolon || char == comma || char == period:
			parsed = append(parsed, sql[idx:idx+1])

		} // end switch

	} // end loop rune range

	if parser.commentDepth > 0 {
		return nil, errors.New("unclosed comment")
	}

	if parser.quoteType != stringTypeNone {
		return nil, errors.New("unclosed quote")
	}

	if parser.skipTo > 0 {
		parsed = append(parsed, strings.ToUpper(sql[parser.stringStartChar:]))
	}

	return parsed, nil

}

func checkDollarTag(parser *parserData, start int) bool {
	if parser.dollarQuote = captureTag(string(parser.chars)[start:]); parser.dollarQuote != nil {
		parser.skipTo = start + parser.dollarQuote.end
		return true
	}
	return false
}

var re = regexp.MustCompile(`(?m)\$([a-zA-Z_0-9]*?)\$`)

type tagResult struct {
	tag      string
	contents string
	end      int
}

// returns a *tagResult if a valid matching tag string was found; including $$$$ (empty tag, empty contents)
func captureTag(str string) *tagResult {

	rslt := re.FindAllStringSubmatchIndex(str, -1)

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
			return &tagResult{
				tag:      tagStr,
				contents: str[rslt[0][1]:rslt[i][0]],
				end:      rslt[i][3],
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
