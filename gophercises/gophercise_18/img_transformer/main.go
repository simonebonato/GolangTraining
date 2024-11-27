package main

import (
	"fmt"
	"img_transformer/primitive"
	"io"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo"
)

func main() {

	// create a static directory used to serve files to the server
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		os.Mkdir("static", 0755)
	}

	e := echo.New()

	// html to submit a new picture using a form and a POST request
	e.GET("/", handleInitialForm)
	e.POST("/upload", handleUpload)

	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":3001"))
}

type TemplateFormData struct {
	Script    string
	ModeNames map[primitive.Mode]string
	NShapes   int
}

func handleInitialForm(c echo.Context) error {
	errorMessage := c.QueryParam("error")
	script := ""
	if errorMessage != "" {
		// Escape the error message for safe inclusion in JavaScript
		escapedMessage := template.JSEscapeString(errorMessage)
		script = fmt.Sprintf(`<script>alert("%s");</script>`, escapedMessage)
	}

	// create the data to pass to the template
	data := TemplateFormData{
		Script:    script,
		ModeNames: primitive.ModeNames,
		NShapes:   100, // default value for the shapes
	}

	tmpl, err := template.ParseFiles("html/form.tmpl")
	if err != nil {
		panic(err)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return tmpl.Execute(c.Response().Writer, data)

}

func handleUpload(c echo.Context) error {
	// check for the form values needed for the transform
	// TODO: use this value to actually change the mode
	mode := c.FormValue("mode")
	mode_int, _ := strconv.Atoi(mode)
	mode_mode := primitive.Mode(mode_int)

	// check that n_shapes is an actual number, not text
	n_shapes_str := c.FormValue("N")
	N, err := strconv.Atoi(n_shapes_str)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/?error=Please+set+N+shapes+as+integer")
	}

	// load the image from the form, if there is one
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println("Error loading the image!")
		return c.Redirect(http.StatusSeeOther, "/?error=Please+select+an+image+to+upload")
	}

	// open it into a file
	f, err := file.Open()
	if err != nil {
		fmt.Println("Error opening the src!")
		return err
	}
	defer f.Close()

	// transform the image and copy the result
	out, err := primitive.Transform(
		f, N, primitive.WithMode(mode_mode),
	)
	if err != nil {
		panic(err)
	}

	// create a file in the static folder
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("static/out_%d.png", timestamp)

	// Create destination file in the static folder
	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// copy the image from the io.Reader to the destination file
	io.Copy(dst, out)

	// call the other handler to display the image!
	err = htmlDisplayImg(c, timestamp)
	if err != nil {
		panic(err)
	}

	return nil
}

func htmlDisplayImg(c echo.Context, timestamp int64) error {
	// Construct the image URL
	imgURL := fmt.Sprintf("/static/out_%d.png", timestamp)

	// Prepare the HTML response
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Transformed Image</title>
	</head>
	<body>
		<h1>Transformed Image</h1>
		<img src='%s' alt='Transformed Image'>
		<br>
		<a href="/">Back to home.</a>
	</body>		
	</html>
`, imgURL)

	return c.HTML(http.StatusOK, htmlContent)
}
