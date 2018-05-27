package main

import (
	"golang.org/x/tour/reader"
)

type MyReader struct{}

// TODO: Add a Read([]byte) (int, error) method to MyReader.
func (myReader MyReader) Read(b []byte) (int,error){
	b = b[:cap(b)];
    for i := range b {
        b[i] = 'A'
    }
	return cap(b), nil
}

func main() {
	//probably calles the Read method of MyReader
	reader.Validate(MyReader{})
}