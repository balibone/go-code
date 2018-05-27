package main

import (
	"fmt"
	"math/rand"
)

func main() {
	for index, count := 0, 10; index < count; index++ {
		fmt.Printf("hello number %d\n", rand.Intn(10))
	}
}
