package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	r := strings.NewReader("Hello, Reader!")

	b := make([]byte, 8)
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("byte slice contents, from index 0 to index %v is %q\n", n-1, b[:n])
		if err == io.EOF {
			break
		}
	}
}