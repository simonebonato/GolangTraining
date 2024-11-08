package cyoa

var defaultHanderTmpl string = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Choose Your Own Adventure</title>
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
            <p>{{.}}</p>
        {{end}}
        <ul>
            {{range .Options}}
                <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
            {{end}}
        </ul>
        <br><br>
        <a href="/">Back to the beginning.</a>
    </body>
</html>`

// we use the "html render golang" package to render the story to the webpage 
// type Chapter struct {
// 	Title	    string   `json:"title"`
// 	Paragraphs  []string `json:"story"`
// 	Options 	[]Option `json:"options"`
// }

// type Option struct {
// 	Text string `json:"text"`
// 	Chapter string `json:"arc"`
// }


var StoryTmpl string = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Choose Your Own Adventure</title>
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
            <p>{{.}}</p>
        {{end}}
        <ul>
            {{range .Options}}
                <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
            {{end}}
        </ul>
        <br><br>
        <a href="/story">Back to the beginning.</a>
    </body>
</html>`