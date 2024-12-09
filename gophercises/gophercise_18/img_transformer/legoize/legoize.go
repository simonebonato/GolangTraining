package legoize

import (
	"fmt"
	"img_transformer/utils"
	"io"
	"os"
	"os/exec"
	"strconv"
)

func WithSize(size int) func() []string {
	return func() []string {
		return []string{"-size", strconv.Itoa(int(size))}
	}
}

func Transform(image io.Reader, legoColors int, file_extension string, opts ...func() []string) (io.Reader, error) {
	in, out, err := utils.CreateInputOutputForTransform(file_extension)
	if err != nil {
		panic(err)
	}

	// to make sure that the file is deleted at the end
	defer os.Remove(in.Name())
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
	_, err = ExecLego(in.Name(), out.Name(), legoColors, opts...)
	if err != nil {
		return nil, err
	}
	// for flushing
	out.Close()
	out, buffer, err := utils.ReadTransformOutput(out)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// we return the buffer, because IT CAN become a Reader
	// and we do not return the output file because it is deleted after the program runs
	return buffer, nil
}

func ExecLego(input_file string, output_file string, lego_colors int, opts ...func() []string) (string, error) {
	cmdString := []string{"-in", input_file, "-out", output_file, "-colors", strconv.Itoa(lego_colors)}

	for _, opt := range opts {
		cmdString = append(cmdString, opt()...)
	}

	cmd := exec.Command("legoizer", cmdString...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(b))
		panic(err)
	}
	return string(b), err
}
