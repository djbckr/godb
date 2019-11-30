package token

import (
	"fmt"
	"testing"
)

func check(t *testing.T, p []string, values ...string) {
	for k, v := range p {
		if values[k] != v {
			t.Error(fmt.Printf("%v != %v", values[k], v))
		}
	}
}

func TestTokenize(t *testing.T) {

	p, e := Tokenize(`select*from dual`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, "SELECT", "*", "FROM", "DUAL")

	p, e = Tokenize(`/* this is a comment */select[issue]from"something"where ix='fubar'`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, "SELECT", "issue", "FROM", "something", "WHERE", "IX", "=", "'fubar'")

	p, e = Tokenize(`  -- single line comment
alter table "schema"."tbl" add column (colname number, colx string) `)

	if e != nil {
		t.Error(e)
	}

	check(t, p, "ALTER", "TABLE", "schema", ".", "tbl", "ADD", "COLUMN", "(", "COLNAME", "NUMBER", ",", "COLX", "STRING", ")")

	p, e = Tokenize(`with x as (select * from dual)
select * from x`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, "WITH", "X", "AS", "(", "SELECT", "*", "FROM", "DUAL", ")", "SELECT", "*", "FROM", "X")

	p, e = Tokenize(`$$
begin do something end;
$$`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, "\nbegin do something end;\n")

	p, e = Tokenize("`This is a test string`")

	if e != nil {
		t.Error(e)
	}

	check(t, p, "This is a test string")

	p, e = Tokenize(`$$$$`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, "")

	p, e = Tokenize(`$tag$xyz$tag$`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, "xyz")


	p, e = Tokenize(`create table [foo$$]
(col1 number, col2 varchar(33, 23))
`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, "CREATE", "TABLE", "foo$$", "(", "COL1", "NUMBER", ",", "COL2", "VARCHAR", "(", "33", ",", "23", ")", ")")

	p, e = Tokenize(`$$begin do something end;$$`)

	if e != nil {
		t.Error(e)
	}

	check(t, p, "begin do something end;")

}
