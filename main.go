package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anshujalan/url-shortener/urlshort"
	"github.com/boltdb/bolt"
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

func buildDBHandler(fallback http.HandlerFunc) (http.HandlerFunc, error) {
	prepareDB()

	return urlshort.DBHandler([]byte("main"), fallback)
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

	dbHander, err := buildDBHandler(jsonHander)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHander)
}

func prepareDB() {
	db, err := bolt.Open("path.db", 0777, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("main"))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte("/bolt"), []byte("https://github.com/boltdb/bolt"))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte("/intro"), []byte("https://npf.io/2014/07/intro-to-boltdb-painless-performant-persistence/"))
		if err != nil {
			return err
		}

		return nil
	})
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
