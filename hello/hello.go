package main

import (
	"io"
	"os"
	"strings"

	"fmt"
	"image"
)

type rt struct {
	r io.Reader
}

func (v rt) Read(b []byte) (n int, err error) {
	fmt.Println("test")
	cnt, ok := v.r.Read(b)
	return cnt, ok
}

func main() {
	s := strings.NewReader("azAZLbh penpxrq gur pbqr!")
	r := rt{s}
	io.Copy(os.Stdout, &r)
	image.NewRGBA()
}
