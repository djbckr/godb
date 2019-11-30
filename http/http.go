package http

import (
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

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/login", login)
}
