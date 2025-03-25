package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	file_name := "to_download.txt"
	MAX_PARALLEL := 10

	f, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var wg sync.WaitGroup
	var sem = make(chan int, MAX_PARALLEL)

	for scanner.Scan() {
		// obtaine the link and the dest of the file
		link := scanner.Text()
		scanner.Scan()
		dest := scanner.Text()

		// check if the destination exists, in that case skip it
		if _, dest_exists := os.Stat(dest); dest_exists == nil {
			fmt.Println("File skipped because it exists already.")
			continue
		}

		// check if the parent dir exists, if not create it
		parent_path := filepath.Dir(dest)
		_, path_exists := os.Stat(parent_path)
		if path_exists != nil {
			os.MkdirAll(parent_path, 0755)
		}

		// download the file!
		sem <- 1
		wg.Add(1)

		go func() {
			defer wg.Done()

			fmt.Printf("\nDownloading file to dest: %s", dest)

			// get the file from the link
			resp, err := http.Get(link)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			// read the body of the response and write to file
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			os.WriteFile(dest, body, 0755)
			<-sem
		}()

	}

	wg.Wait()

}
