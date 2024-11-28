package primitive

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
)

// Defines the shapes used when transforming images
type Mode int

// Modes supported by the primitive package
const (
	ModeTriangle Mode = iota
	ModeCombo
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

var ModeNames = map[Mode]string{
	ModeTriangle:       "Triangle",
	ModeCombo:          "Combo",
	ModeRect:           "Rect",
	ModeEllipse:        "Ellipse",
	ModeCircle:         "Circle",
	ModeRotatedRect:    "Rotated Rect",
	ModeBeziers:        "Beziers",
	ModeRotatedEllipse: "Rotated Ellipse",
	ModePolygon:        "Polygon",
}

// Option for the transform function for the mode you want to use,
// by default mode triangle will be used
func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", strconv.Itoa(int(mode))}
	}
}

// what I lack:
// - what is an io.Reader and how to use it
// - what is a buffer instead
// - what is a tempFile for
// - why are we using these things and not simply loading an image as in Python?

// Transform will take the provided image as an io.Reader, apply the primitive transformation,
// and return an io.Reader with the corresponding transformed image
func Transform(image io.Reader, numShapes int, file_extension string, opts ...func() []string) (io.Reader, error) {
	// create a tmp file for the input image
	in, err := os.CreateTemp("", fmt.Sprintf("in_*%s", file_extension))
	if err != nil {
		return nil, err
	}
	// to make sure that the file is deleted at the end
	defer os.Remove(in.Name())

	// do the same for the output file
	out, err := os.CreateTemp("", fmt.Sprintf("out_*%s", file_extension))
	if err != nil {
		return nil, err
	}
	defer os.Remove(out.Name())

	// read image into input files
	_, err = io.Copy(in, image)
	if err != nil {
		return nil, err
	}
	// we need to close the file for flushing
	// layman terms: the system will finish writing all we need to write
	// to the file, and we are sure in the next operation the file is good
	in.Close()

	// run primitive w/ -i in.Name() -o out.Name()
	_, err = ExecPrimitive(in.Name(), out.Name(), numShapes, opts...)
	if err != nil {
		return nil, err
	}
	// for flushing
	out.Close()

	// Reopen 'out' for reading
	out, err = os.Open(out.Name())
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// Read from 'outputFile' into the buffer
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, out)
	if err != nil {
		return nil, err
	}

	// we return the buffer, because IT CAN become a Reader
	// and we do not return the output file because it is deleted after the program runs
	return buffer, nil
}

func ExecPrimitive(input_file string, output_file string, num_shapes int, opts ...func() []string) (string, error) {
	cmdString := []string{"-i", input_file, "-o", output_file, "-n", strconv.Itoa(num_shapes)}

	for _, opt := range opts {
		cmdString = append(cmdString, opt()...)
	}

	cmd := exec.Command("primitive", cmdString...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(b))
		panic(err)
	}
	return string(b), err
}
