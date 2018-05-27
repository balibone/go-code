package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (rot13 rot13Reader) Read(b []byte) (n int, err error){
	n,err = rot13.r.Read(b)
	for i := 0 ; i < n ; i++ {
		//implement ROT 13
		// if byte read (character read) is between A to Z
		if 'A' <= b[i] && b[i] <= 'Z' {
			b[i] = (b[i] - 'A' + 13) % 26 + 'A'
		}
		// if byte read (character read) is between a to z
		if 'a' <= b[i] && b[i] <= 'z' {
			b[i] = (b[i] - 'a' + 13) % 26 + 'a'
		}
	}
	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}//initialises a rot13Reader containing an io.Reader pointed to that random string above
	io.Copy(os.Stdout, &r)
}
