package main

import(
	"fmt"
	"math/rand"
)

func main() {
	count int = 10
	for index := 0; index < count; index++ {
		fmt.Println("hello number", rand.Int)
	}
}
