package cyoa

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

// ------------- types ---------------
type Story map[string]Chapter

type Chapter struct {
	Title	    string   `json:"title"`
	Paragraphs  []string `json:"story"`
	Options 	[]Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Chapter string `json:"arc"`
}

// ------------- utils ----------------

func JSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// ------------- http funcs ---------------- 

// we can create a type that can modify a handler, since it gets a pointer to it
// this is the base of functional programming
type handlerOption func (*handler)

// function closure, to make the code cleaner and easier for the end user
func WithTemplate (t *template.Template) handlerOption {
	return func (h *handler) {
		h.t = t
	}
}

// technically we can create other functions with some optional argument for the arguments
// eg. func WithDatabase (*db) handlerOption {}

// in the ...handlerOption we can pass in multiple options!
func NewHandler (story Story, opts ...handlerOption) http.Handler {
	// the must is to control that the template is prod ready, correct
	tpl := template.Must(template.New("HTMLStoryTemplate").Parse(defaultHanderTmpl))
	h := handler{story, tpl}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

type handler struct {
	story Story
	t *template.Template
}

func (h handler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	
	url := strings.TrimSpace(r.URL.Path)

	var render_chapter string
	if url == "/" {
		render_chapter = "intro"
	} else {
		render_chapter = url[1:]
	}

	if chapter, ok := h.story[render_chapter]; ok {
		err := h.t.Execute(w, chapter)

		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong dude...", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found...", http.StatusNotFound)
}


