package dml

import (
	"fmt"
	"github.com/djbckr/godb/sql/token"
	"testing"
)

func TestProcessSelect(t *testing.T) {
	tokens, e := token.Tokenize(`select * from dual where 1 = 1.0`)

	fmt.Println("-------------")

	if e != nil {
		t.Error(e)
	}

	ProcessSelect(tokens)

	tokens, e = token.Tokenize(`select * from dual where x = q'[bob]'`)

	fmt.Println("-------------")

	if e != nil {
		t.Error(e)
	}

	ProcessSelect(tokens)


}
