package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateInputOutputForTransform(file_extension string) (*os.File, *os.File, error) {
	// create a tmp file for the input image
	in, err := os.CreateTemp("", fmt.Sprintf("in_*%s", file_extension))
	if err != nil {
		return nil, nil, err
	}

	// do the same for the output file
	out, err := os.CreateTemp("", fmt.Sprintf("out_*%s", file_extension))
	if err != nil {
		return nil, nil, err
	}

	return in, out, nil
}

func ReadTransformOutput(out *os.File) (*os.File, *bytes.Buffer, error) {
	// Reopen 'out' for reading
	out, err := os.Open(out.Name())
	if err != nil {
		return nil, nil, err
	}

	// Read from 'outputFile' into the buffer
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, out)
	if err != nil {
		return nil, nil, err
	}

	return out, buffer, err
}

// Function to clean the static folder
func CleanStaticFolder(folder string) error {
	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.Remove(filepath.Join(folder, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to remove file %s: %w", file.Name(), err)
		}
	}

	fmt.Println("Static folder cleaned up.")
	return nil
}
