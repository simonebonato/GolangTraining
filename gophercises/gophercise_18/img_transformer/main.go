package main

import (
	"fmt"
	"img_transformer/primitive"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/labstack/echo"
)

var imagePath string

func main() {

	// create a static directory used to serve files to the server
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		os.Mkdir("static", 0755)
	}

	e := echo.New()

	// html to submit a new picture using a form and a POST request
	e.GET("/", handleInitialForm)
	e.POST("/upload", handleUpload)
	// TODO: add something like this for handling the transform and remove the global variable
	// e.POST("/transform/:primitive:filename", handlePrimitiveTransform)
	e.GET("/display/:filename", handleHtmlDisplayImg)

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

	imagePath = ""

	// load the image from the form, if there is one
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println("Error loading the image!")
		return c.Redirect(http.StatusSeeOther, "/?error=Please+select+an+image+to+upload")
	}
	file_extension := filepath.Ext(file.Filename)

	// open it into a file
	f, err := file.Open()
	if err != nil {
		fmt.Println("Error opening the src!")
		return err
	}
	defer f.Close()

	// save the image locally and store the path in imagePath
	img_filename := createStaticTimestampFilename(file_extension, "")
	tmp_file, err := os.Create(img_filename)
	if err != nil {
		return err
	}
	defer tmp_file.Close()

	_, err = io.Copy(tmp_file, f)
	if err != nil {
		return err
	}

	// now this can technically be used everywhere in these functions since it is a global variable
	imagePath = img_filename

	// to check which transform has been triggered
	transformType := c.FormValue("transform")

	if transformType == "primitive" {
		err = handlePrimitiveTransform(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func handlePrimitiveTransform(c echo.Context) error {
	// check if the image path is set
	if imagePath == "" {
		fmt.Println("No image loaded!")
		return c.Redirect(http.StatusSeeOther, "/?error=Please+select+an+image+to+upload")
	}

	// read the image into a io.Reader for the transform method
	f, err := os.Open(imagePath)
	file_extension := filepath.Ext(imagePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// check for the form values needed for the transform
	mode := c.FormValue("mode")
	mode_int, _ := strconv.Atoi(mode)
	mode_mode := primitive.Mode(mode_int)

	// check that n_shapes is an actual number, not text
	n_shapes_str := c.FormValue("N")
	N, err := strconv.Atoi(n_shapes_str)
	if err != nil {
		fmt.Println("Error loading the image!")
		return c.Redirect(http.StatusSeeOther, "/?error=Please+set+N+shapes+as+integer")
	}

	// transform the image and copy the result
	out, err := primitive.Transform(
		f, N, file_extension, primitive.WithMode(mode_mode),
	)
	if err != nil {
		panic(err)
	}

	// now save the image, then display it create a file in the static folder
	// TODO: maybe turn this into a function that can be used by both the primitive and Legoize transforms
	filename := createStaticTimestampFilename(file_extension, "primitive_")
	out_file, err := os.Create(filename)
	if err != nil {
		return err
	}
	io.Copy(out_file, out)

	// call the other handler to display the image!
	relativeFilename := strings.TrimPrefix(filename, "static/")
	err = c.Redirect(http.StatusSeeOther, fmt.Sprintf("/display/%s", relativeFilename))
	if err != nil {
		panic(err)
	}

	return nil
}

func createStaticTimestampFilename(file_extension string, prefix string) string {
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("static/%sout_%d%s", prefix, timestamp, file_extension)
	return filename
}

func handleHtmlDisplayImg(c echo.Context) error {

	filename := c.Param("filename")

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
		<img src='/static/%s' alt='Transformed Image'>
		<br>
		<a href="/">Back to home.</a>
	</body>		
	</html>
`, filename)

	return c.HTML(http.StatusOK, htmlContent)
}
