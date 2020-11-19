package main

import (
	"fmt"
	"net/http"

	"github.com/anshujalan/url-shortener/urlshort"
)

func buildMapHandler(fallback http.Handler) http.HandlerFunc {
	pathsToUrls := map[string]string{
		"/urlshort-doc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-doc":     "https://pkg.go.dev/gopkg.in/yaml.v2",
	}
	return urlshort.MapHandler(pathsToUrls, fallback)
}

func buildYAMLHander(fallback http.Handler) (http.HandlerFunc, error) {
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	return urlshort.YAMLHandler([]byte(yaml), fallback)
}

func buildJSONHandler(fallback http.HandlerFunc) (http.HandlerFunc, error) {
	jsonData := `[
		{
			"path": "/golang",
			"url": "https://en.wikipedia.org/wiki/Go_(programming_language)"
		},
		{
			"path": "/json",
			"url": "https://golang.org/pkg/encoding/json/"
		}
	]`

	return urlshort.JSONHandler([]byte(jsonData), fallback)
}

func main() {
	mux := defaultMux()

	mapHandler := buildMapHandler(mux)

	yamlHandler, err := buildYAMLHander(mapHandler)
	if err != nil {
		panic(err)
	}

	jsonHander, err := buildJSONHandler(yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHander)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
