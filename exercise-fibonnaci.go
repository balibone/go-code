package main

import "fmt"

// fibonacci is a function that returns
// a closure function that returns an int.
func fibonacci() func() int {
	f1 := -1 //start f1 at -1 so output can print 0,1,1
	f2 := 0
	return func() int{
		if(f1 == -1 && f2 == 0){//print 0
			f1++
			return 0
		}else if(f1 == 0 && f2 == 0){//print 1st 1
			f2++
			return 1
		}
		sum := f1+f2
		f1 = f2
		f2 = sum
		return sum
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
