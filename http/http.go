package http

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

func root(rsp http.ResponseWriter, req *http.Request)  {

}

func admin(rsp http.ResponseWriter, req *http.Request)  {

}

func login(rsp http.ResponseWriter, req *http.Request)  {
	rsp.Header().Add("Authorization", "whoopdedoo");
	_, _ = io.WriteString(rsp, "done")
}

func sql(rsp http.ResponseWriter, req *http.Request) {
	var decoder *Decoder
	switch req.Header.Get("Content-Type") {
	case "application/json":
		decoder = json.NewDecoder(req.Body)
	case "application/xml":
		decoder = xml.NewDecoder(req.Body)

	}
	req.Header.Get("Content-Type") // body
	// req.Header.Get("Accept")
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/login", login)
	http.HandleFunc("/sql", sql)
}
