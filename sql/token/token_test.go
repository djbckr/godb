package token

import (
	"fmt"
	"math/big"
	"testing"
)

func check(t *testing.T, p Tokens, values ...Token) {
	for k, v := range p {
		var test bool
		switch vv := values[k].Value.(type) {
		case string:
			test = vv != v.Value
		case big.Float:
			test = vv.Cmp(v.Value.(*big.Float)) == 0
		}

		if test ||
			values[k].TokenType != v.TokenType{
			t.Error(fmt.Printf("%v != %v\n", values[k], v))
		}
	}
}

func TestTokenize(t *testing.T) {

	p, e := Tokenize(`select*from dual`)

	if e != nil {
		t.Error(e)
	}

	check(t, p,
		Token{"SELECT", TypeToken},
		Token{"*", TypePunctuation},
		Token{"FROM", TypeToken},
		Token{"DUAL", TypeToken})

	p, e = Tokenize(`/* this is a comment */select[issue]from"something"where ix='fubar'`)

	if e != nil {
		t.Error(e)
	}

	check(t, p,
		Token{" this is a comment ", TypeComment},
		Token{"SELECT", TypeToken},
		Token{"issue", TypeToken},
		Token{"FROM", TypeToken},
		Token{"something", TypeToken},
		Token{"WHERE", TypeToken},
		Token{"IX", TypeToken},
		Token{"=", TypePunctuation},
		Token{"fubar", TypeString})

	p, e = Tokenize(`/* this is a comment */select[issue]from"something"where ix=1e7`)

	if e != nil {
		t.Error(e)
	}

	n, _, _ := big.NewFloat(0).Parse("1e7", 0)
	check(t, p,
		Token{" this is a comment ", TypeComment},
		Token{"SELECT", TypeToken},
		Token{"issue", TypeToken},
		Token{"FROM", TypeToken},
		Token{"something", TypeToken},
		Token{"WHERE", TypeToken},
		Token{"IX", TypeToken},
		Token{"=", TypePunctuation},
		Token{n, TypeNumber})

	p, e = Tokenize(`  -- single line comment
alter table "schema"."tbål" add column (colname Number, colx string) `)

	if e != nil {
		t.Error(e)
	}

	check(t, p,
		Token{" single line comment", TypeComment},
		Token{"ALTER", TypeToken},
		Token{"TABLE", TypeToken},
		Token{"schema", TypeToken},
		Token{".", TypePunctuation},
		Token{"tbål", TypeToken},
		Token{"ADD", TypeToken},
		Token{"COLUMN", TypeToken},
		Token{"(", TypePunctuation},
		Token{"COLNAME", TypeToken},
		Token{"NUMBER", TypeToken},
		Token{",", TypePunctuation},
		Token{"COLX", TypeToken},
		Token{"STRING", TypeToken},
		Token{")", TypePunctuation}	)

	p, e = Tokenize(`with x as (select * from dual)
select * from x`)

	if e != nil {
		t.Error(e)
	}

	check(t, p,
		Token{"WITH", TypeToken},
		Token{"X", TypeToken},
		Token{"AS", TypeToken},
		Token{"(", TypePunctuation},
		Token{"SELECT", TypeToken},
		Token{"*", TypePunctuation},
		Token{"FROM", TypeToken},
		Token{"DUAL", TypeToken},
		Token{")", TypePunctuation},
		Token{"SELECT", TypeToken},
		Token{"*", TypePunctuation},
		Token{"FROM", TypeToken},
		Token{"X", TypeToken},
		)

	p, e = Tokenize(`$$
begin do something end;
$$`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, Token{"\nbegin do something end;\n", TypeString})

	p, e = Tokenize("`This is a test string`")

	if e == nil || e.Error() != "Invalid SQL text found near: `This is a test string`" {
		t.Error(e)
	}


	p, e = Tokenize(`$$$$`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, Token{"", TypeString})

	p, e = Tokenize(`$tag$xyz$tag$`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, Token{"xyz", TypeString})

	p, e = Tokenize(`$⌘$†π¬˚$⌘$`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, Token{"†π¬˚", TypeString})


	p, e = Tokenize(`create table [foo$$]
(col1 Number, col2 varchar(33, 23))
`)

	if e != nil {
		t.Error(e)
	}

	n1, _, _ := big.NewFloat(0).Parse("33", 0)
	n2, _, _ := big.NewFloat(0).Parse("23", 0)
	check(t, p,
		Token{"CREATE",TypeToken},
		Token{"TABLE",TypeToken},
		Token{"foo$$",TypeToken},
		Token{"(",TypePunctuation},
		Token{"COL1",TypeToken},
		Token{"NUMBER",TypeToken},
		Token{",",TypePunctuation},
		Token{"COL2",TypeToken},
		Token{"VARCHAR",TypeToken},
		Token{"(",TypePunctuation},
		Token{n1,TypeNumber},
		Token{",",TypePunctuation},
		Token{n2,TypeNumber},
		Token{")",TypePunctuation},
		Token{")",TypePunctuation}	)

	p, e = Tokenize(`$$begin do something end;$$`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, Token{"begin do something end;", TypeString})

	p, e = Tokenize(`select * from something where x = q'[mary's horse]' foo bar`)

	if e != nil {
		t.Error(e)
	}

	check(t, p,
		Token{"SELECT", TypeToken},
		Token{"*", TypePunctuation},
		Token{"FROM", TypeToken},
		Token{"SOMETHING", TypeToken},
		Token{"WHERE", TypeToken},
		Token{"X", TypeToken},
		Token{"=", TypePunctuation},
		Token{"mary's horse", TypeString},
		Token{"FOO", TypeToken},
		Token{"BAR", TypeToken},
		)

	p, e = Tokenize(`select * from something where x = q'[mary's horse]'`)

	if e != nil {
		t.Error(e)
	}

	check(t, p,
		Token{"SELECT", TypeToken},
		Token{"*", TypePunctuation},
		Token{"FROM", TypeToken},
		Token{"SOMETHING", TypeToken},
		Token{"WHERE", TypeToken},
		Token{"X", TypeToken},
		Token{"=", TypePunctuation},
		Token{"mary's horse", TypeString},
	)

	p, e = Tokenize(`q'[]'`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, Token{"", TypeString})

}
