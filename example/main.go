//go:generate gohtml tpl

package main

import (
	"net/http"

	"github.com/Bluek404/gohtml/example/tpl"
)

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write(tpl.Index(r.URL.Path[1:]))
}
