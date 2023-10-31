package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

func main() {
	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":1234", yamlHandler))

}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
	return nil
}
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedPathUrls, err := parseUrls(yamlBytes)
	if err != nil {
		return nil, err
	}
	mappedUrls := mapBuilder(parsedPathUrls)
	return MapHandler(mappedUrls, fallback), err
	// pathUrls, err := parseUrls(yamlBytes)
	// if err != nil {
	// 	return nil, err
	// }
	// pathsToUrls := mapBuilder(pathUrls)
	// return MapHandler(pathsToUrls, fallback), nil
}

func mapBuilder(pathUrls []pathUrl) map[string]string {
	mapOfUrls := make(map[string]string)

	for _, val := range pathUrls {
		mapOfUrls[val.Path] = val.URL
	}

	return mapOfUrls
}

func parseUrls(data []byte) ([]pathUrl, error) {
	var parsedPathUrls []pathUrl
	err := yaml.Unmarshal(data, &parsedPathUrls)
	if err != nil {
		return nil, err
	}

	return parsedPathUrls, err
}

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
