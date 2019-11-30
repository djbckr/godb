package main

import (
	_ "github.com/djbckr/godb/http"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServeTLS(":9422", "/Users/dbecker/.ssh/alchemy.crt", "/Users/dbecker/.ssh/alchemy.key", nil))
}
