package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
)

func main() {

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
	e.GET("/display/:filename", handleDisplayImage)
	e.Logger.Fatal(e.Start(":3001"))
}

func handleUpload(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println("Error loading the image!")
		return err
	}

	src, err := file.Open()
	if err != nil {
		fmt.Println("Error opening the src!")
		return err
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join("uploads", file.Filename))
	if err != nil {
		fmt.Println("Error finding the dst!")
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println("Error copying!")
		return err
	}

	displayURL := fmt.Sprintf("/display/%s", file.Filename)
	return c.HTML(
		http.StatusOK, 
		fmt.Sprintf("<p>File %s uploaded successfully. <a href='%s'>View image</a></p>", file.Filename, displayURL))
}

func handleDisplayImage(c echo.Context) error {
	filename := c.Param("filename")
    filepath := filepath.Join("uploads", filename)

    return c.File(filepath)
}