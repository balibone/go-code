package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	t, z := 0., 1.
    for {
		//z will keep changing to be more accurate, and t will keep taking the old value of z
        z, t = z - (z*z-x)/(2*z), z
		//if difference between t and z becomes minutely small, correct square root has been reached. 
        if math.Abs(t-z) < 1e-8 {
            break
        }
    }
    return z
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(2) == math.Sqrt(2))
}