package sql

import "testing"

func TestParse(t *testing.T) {

	t.Log("-----------")

	p, e := textParse(`select*from dual`)

	if e != nil {
		t.Error(e)
	}

	for _, s := range p {
		t.Log(">", s, "<")
	}

	t.Log("-----------")

	p, e = textParse(`/* this is a comment */select[issue]from"something"where ix='fubar'`)

	if e != nil {
		t.Error(e)
	}

	for _, s := range p {
		t.Log(">", s, "<")
	}

	t.Log("-----------")

	p, e = textParse(`  -- single line comment
alter table "schema"."tbl" add column (colname number, colx string) `)

	if e != nil {
		t.Error(e)
	}

	for _, s := range p {
		t.Log(">", s, "<")
	}

	t.Log("-----------")

	p, e = textParse(`with x as (select * from dual)
select * from x`)

	if e != nil {
		t.Error(e)
	}

	for _, s := range p {
		t.Log(">", s, "<")
	}

	t.Log("-----------")

	p, e = textParse(`$$
begin do something end;
$$`)

	if e != nil {
		t.Error(e)
	}

	for _, s := range p {
		t.Log(">", s, "<")
	}

	t.Log("-----------")

	p, e = textParse("`This is a test string`")

	if e != nil {
		t.Error(e)
	}

	for _, s := range p {
		t.Log(">", s, "<")
	}

	t.Log("-----------")

	p, e = textParse(`$$$$`)

	if e != nil {
		t.Error(e)
	}

	for _, s := range p {
		t.Log(">", s, "<")
	}

	t.Log("-----------")

	p, e = textParse(`create table [foo$$]
(col1 number, col2 varchar(33, 23))
`)

	if e != nil {
		t.Error(e)
	}

	for _, s := range p {
		t.Log(">", s, "<")
	}

	t.Log("-----------")

	p, e = textParse(`$$begin do something end;$$`)

	if e != nil {
		t.Error(e)
	}

	for _, s := range p {
		t.Log(">", s, "<")
	}


}
