package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		des, ok := pathsToUrls[path]

		if ok {
			http.Redirect(w, r, des, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	out, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMapping := buildMapping(out)

	return MapHandler(pathMapping, fallback), nil
}

//JSONHandler builds the http.Handlerfunc for handling json based path:url storage
func JSONHandler(jsn []byte, fallback http.HandlerFunc) (http.HandlerFunc, error) {
	out, err := parseJSON(jsn)
	if err != nil {
		return nil, err
	}
	pathMapping := buildMapping(out)

	return MapHandler(pathMapping, fallback), nil
}

//DBHandler builds the http.HandlerFunc from local BoltDB storage
func DBHandler(bkt []byte, fallback http.HandlerFunc) (http.HandlerFunc, error) {
	db, err := bolt.Open("path.db", 0777, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var val []byte
		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(bkt)

			if bucket == nil {
				fallback.ServeHTTP(w, r)
				return nil
			}
			val = bucket.Get([]byte(r.URL.Path))

			return nil
		})

		if err != nil {
			fallback.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, string(val), http.StatusFound)
	}, nil
}

type pathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func buildMapping(out []pathURL) map[string]string {
	pathMapping := make(map[string]string)
	for _, pu := range out {
		pathMapping[pu.Path] = pu.URL
	}

	return pathMapping
}

func parseYAML(yml []byte) ([]pathURL, error) {
	out := []pathURL{}
	err := yaml.Unmarshal(yml, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func parseJSON(jsn []byte) ([]pathURL, error) {
	out := []pathURL{}
	err := json.Unmarshal(jsn, &out)

	if err != nil {
		return nil, err
	}

	return out, nil
}
