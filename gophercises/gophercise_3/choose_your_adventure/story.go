package cyoa

import (
	"encoding/json"
	"io"
	"net/http"
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

func NewHandler (story Story) http.Handler {
	return handler{story:story}
}

type handler struct {
	story Story
}

func (h handler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	// the must is to control that the template is prod ready, correct
	tpl := template.Must(template.New("HTMLStoryTemplate").Parse(defaultHanderTmpl))

	url := r.URL.Path

	var render_chapter string
	if url == "/" {
		render_chapter = "intro"
	} else {
		render_chapter = url[1:]
	}

	err := tpl.Execute(w, h.story[render_chapter])
	if err != nil {
		panic(err)
	}
}