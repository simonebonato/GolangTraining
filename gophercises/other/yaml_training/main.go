package main 

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// since the structure of the yaml is known, I can simply define a struct with the corresponding elements
type PathURL struct {
	Path string `yaml:"path"`
	URL string `yaml:"url"`
}

func main() {
		yaml_file := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
		// then create a slice that will contain all the pairs in the slice
		var pathURLS []PathURL

		// finally we can pass the bytes of the string and the pointer to the define slice
		// this way it can be populated with the content of the yml after calling the Unmarshal method
		err := yaml.Unmarshal([]byte(yaml_file), &pathURLS)
		if err != nil {
			fmt.Println("Error unmarshaling the YAML file")
			return
		}
		
		// they are saved in a slice, with the attributes we defined in the struct
		for _, p := range pathURLS {
			fmt.Println(p.Path, p.URL)
		}
}