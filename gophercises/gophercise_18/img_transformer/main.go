package main

import (
	"fmt"
	"img_transformer/legoize"
	"img_transformer/primitive"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	e.GET("/display/:filename", handleHtmlDisplayImg)

	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":3001"))
}

type TemplateFormData struct {
	Script        string
	ModeNames     map[primitive.Mode]string
	NShapes       int
	LastImagePath string
}

func handleInitialForm(c echo.Context) error {
	errorMessage := c.QueryParam("error")
	script := ""
	if errorMessage != "" {
		// Escape the error message for safe inclusion in JavaScript
		escapedMessage := template.JSEscapeString(errorMessage)
		script = fmt.Sprintf(`<script>alert("%s");</script>`, escapedMessage)
	}

	// get the cookie for the image path
	lastImagePath := ""
	cookie, err := c.Cookie("imagePath")
	if err == nil {
		lastImagePath = cookie.Value
	}

	// create the data to pass to the template
	data := TemplateFormData{
		Script:        script,
		ModeNames:     primitive.ModeNames,
		NShapes:       100, // default value for the shapes
		LastImagePath: lastImagePath,
	}

	tmpl, err := template.ParseFiles("html/form.tmpl")
	if err != nil {
		panic(err)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return tmpl.Execute(c.Response().Writer, data)

}

func handleUpload(c echo.Context) error {

	var file_extension string
	var f multipart.File
	var img_filename string

	useLastImage := c.FormValue("useLastImage") == "true"

	if useLastImage {
		existingPath, _ := c.Cookie("imagePath")
		f, err := os.Open(existingPath.Value)
		if err != nil {
			return err
		}
		defer f.Close()
		img_filename = existingPath.Value

	} else {
		// load the image from the form, if there is one
		file, err := c.FormFile("image")
		if err != nil {
			fmt.Println("Error loading the image!")
			return c.Redirect(http.StatusSeeOther, "/?error=Please+select+an+image+to+upload")
		}
		file_extension = filepath.Ext(file.Filename)

		// open it into a file
		f, err = file.Open()
		if err != nil {
			fmt.Println("Error opening the src!")
			return err
		}
		defer f.Close()

		// save the image locally and store the path in imagePath
		img_filename = createStaticTimestampFilename(file_extension, "")
		tmp_file, err := os.Create(img_filename)
		if err != nil {
			return err
		}
		defer tmp_file.Close()

		_, err = io.Copy(tmp_file, f)
		if err != nil {
			fmt.Println("Error copying the image here!!!!")
			return err
		}
	}

	// Store imagePath in the context and in a cookie
	c.Set("imagePath", img_filename)
	writeCookie(c, img_filename)

	// to check which transform has been triggered
	transformType := c.FormValue("transform")

	var err error
	if transformType == "primitive" {
		err = applyPrimitiveTransform(c)
	} else if transformType == "lego" {
		err = applyLegoTransform(c)
	}

	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/?error="+url.QueryEscape(err.Error()))
	}

	return nil
}

func writeCookie(c echo.Context, cookieName string) error {
	cookie := new(http.Cookie)
	cookie.Name = "imagePath"
	cookie.Value = cookieName
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	return nil
}

func getImageFromContext(c echo.Context) (*os.File, string, error) {
	// check if the image path is set
	imagePath, ok := c.Get("imagePath").(string)
	if !ok || imagePath == "" {
		return nil, "", fmt.Errorf("no image loaded")
	}

	// read the image into a io.Reader for the transform method
	f, err := os.Open(imagePath)
	if err != nil {
		return nil, "", fmt.Errorf("no image loaded")
	}
	file_extension := filepath.Ext(imagePath)

	return f, file_extension, nil
}

func saveAndRedirect(c echo.Context, out io.Reader, file_extension string, prefix string) error {
	filename := createStaticTimestampFilename(file_extension, prefix)
	out_file, err := os.Create(filename)
	if err != nil {
		return err
	}
	io.Copy(out_file, out)

	// call the other handler to display the image!
	relativeFilename := strings.TrimPrefix(filename, "static/")
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/display/%s", relativeFilename))
}

type TransformFunc func(input io.Reader, ext string) (io.Reader, error)

func processTransform(c echo.Context, transform TransformFunc, prefix string) error {
	f, ext, err := getImageFromContext(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/?error="+url.QueryEscape(err.Error()))
	}
	defer f.Close()

	out, err := transform(f, ext)
	if err != nil {
		return err
	}

	return saveAndRedirect(c, out, ext, prefix)
}

func applyLegoTransform(c echo.Context) error {
	legoColorsStr := c.FormValue("lego_colors")
	legoSizeStr := c.FormValue("lego_size")

	legoColors, err := strconv.Atoi(legoColorsStr)
	if err != nil { return fmt.Errorf("lego color must be integer")
	}

	legoSize, err := strconv.Atoi(legoSizeStr)
	if err != nil { return fmt.Errorf("lego size must be integer")
	}

	transformFunc := func(input io.Reader, ext string) (io.Reader, error) {
		return legoize.Transform(input, legoColors, ext, legoize.WithSize(legoSize))
	}

	return processTransform(c, transformFunc, "lego")
}

func applyPrimitiveTransform(c echo.Context) error {
	modeStr := c.FormValue("mode")
	nShapesStr := c.FormValue("N")

	modeVal, err := strconv.Atoi(modeStr)
	if err != nil { return fmt.Errorf("select a valid mode")
	}

	N, err := strconv.Atoi(nShapesStr)
	if err != nil { return fmt.Errorf("the number of shapes must be integer")
	}

	transformFunc := func(input io.Reader, ext string) (io.Reader, error) {
		return primitive.Transform(input, N, ext, primitive.WithMode(primitive.Mode(modeVal)))
	}

	return processTransform(c, transformFunc, "primitive")
}

func createStaticTimestampFilename(file_extension string, prefix string) string {
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("static/%s_out_%d%s", prefix, timestamp, file_extension)
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
