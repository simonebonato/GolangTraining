package main

import (
	"fmt"
	"img_transformer/primitive"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
)

func main() {

	// create a static directory used to serve files to the server
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		os.Mkdir("static", 0755)
	}

	// html to submit a new picture
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
	        <form action="/upload" method="post" enctype="multipart/form-data">
	            <input type="file" name="image" accept="image/*">
	            <button type="submit">Upload</button>
	        </form>
	    `)
	})

	e.POST("/upload", handleUpload)
	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":3001"))
}

func handleUpload(c echo.Context) error {
	// load the image from the form
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println("Error loading the image!")
		return err
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
		f, 20,
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

	// Prepare the HTML response
	w := c.Response().Writer
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintln(w, "<html><body>")
	fmt.Fprintln(w, "<h1>Transformed Image</h1>")

	// Construct the image URL
	imgURL := fmt.Sprintf("/static/out_%d.png", timestamp)
	outString := fmt.Sprintf("<img src='%s' alt='Transformed Image'>", imgURL)
	fmt.Fprint(w, outString)

	fmt.Fprintln(w, "</body></html>")

	return nil
}
