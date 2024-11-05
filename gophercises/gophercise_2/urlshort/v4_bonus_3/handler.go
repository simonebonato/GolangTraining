package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

// since the structure of the yaml is known, I can simply define a struct with the corresponding elements
type PathURLYAML struct {
	Path string `yaml:"path"`
	URL string `yaml:"url"`
}

type PathURLJSON struct {
	Path string `json:"path"`
	URL string `json:"url"`
}

type Router struct  {
	pathsToUrls map[string]string
	fallback http.Handler
}

func (router Router) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	val, ok := router.pathsToUrls[path]

	if ok {
		http.Redirect(w, r, val, http.StatusSeeOther)
	} else {
		router.fallback.ServeHTTP(w, r)
	}
}

type DbRouter struct {
	db *gorm.DB
	fallback http.Handler
}

func (dbrouter DbRouter) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	value, exists :=  GetValues(dbrouter.db, path)
	if !exists {
		http.Redirect(w, r, value, http.StatusSeeOther)
	} else {
		dbrouter.fallback.ServeHTTP(w, r)
	}
}

func DbHandler(db *gorm.DB, fallback http.Handler) http.HandlerFunc {
	var dbrouter DbRouter = DbRouter{db : db, fallback: fallback}
	return dbrouter.ServeHTTP
}


// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	var router Router = Router{pathsToUrls: pathsToUrls, fallback: fallback}

	return router.ServeHTTP
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	pathMap, _ := getPathMapYAML(yml)

	// TODO: Implement this...
	return MapHandler(pathMap, fallback), nil
}

func getPathMapYAML(yml []byte) (map[string]string, error) {
	// then create a slice that will contain all the pairs in the slice
	var pathURLS []PathURLYAML

	// finally we can pass the bytes of the string and the pointer to the define slice
	// this way it can be populated with the content of the yml after calling the Unmarshal method
	err := yaml.Unmarshal(yml, &pathURLS)

	if err != nil {
		fmt.Println("Error unmarshaling the YAML file")
		return nil, err
	}
	pathMap := make(map[string]string)
	
	// they are saved in a slice, with the attributes we defined in the struct
	for _, p := range pathURLS {
		pathMap[p.Path] = p.URL
	}

	return pathMap, nil
}


func YAMLFileHandler(yaml_path string, fallback http.Handler) (http.HandlerFunc, error) {
	if yaml_path == "" {
		return fallback.ServeHTTP, nil
	}

	// load the file
	content, err := os.ReadFile(yaml_path)
	if err != nil {
		fmt.Println("Error reading yaml file.", err)
		return nil, err
	}

	pathMap, _ := getPathMapYAML(content)

	return MapHandler(pathMap, fallback), nil

}

// ------------ JSON -----------------

func getPathMapJSON(json_bytes []byte) (map[string]string, error) {

	var paths []PathURLJSON

	err := json.Unmarshal(json_bytes, &paths)

	if err != nil {
		fmt.Println("Error unmarshaling the YAML file")
		return nil, err
	}
	pathMap := make(map[string]string)
	
	// they are saved in a slice, with the attributes we defined in the struct
	for _, p := range paths {
		pathMap[p.Path] = p.URL
	}

	return pathMap, nil
}

func JSONFIleHandler(json_path string, fallback http.Handler) (http.HandlerFunc, error) {
	if json_path == "" {
		return fallback.ServeHTTP, nil
	}

	// load the json file
	content, err := os.ReadFile(json_path)
	if err != nil {
		fmt.Println("Error reading json file.", err)
		return nil, err
	}

	pathMap, _ := getPathMapJSON(content)
	return MapHandler(pathMap, fallback), nil
}