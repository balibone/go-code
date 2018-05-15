package main

func CopyDigits(filename string) []byte{
	b, _ := ioutil.ReadFile(filename)
	b = digitRegexp.Find(b)
	//make empty slice c with length of b.
	c := make([]byte, len(b))
	//truncate b
	b = append(c, b...)
}