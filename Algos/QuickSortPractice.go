package Algos

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	arr := []int{3, 5, 2, 19, 2, 3532, 0, -32}
	quickSort(arr, 0, len(arr)-1)
	fmt.Println("after", arr)
}

func quickSort(arr []int, left, right int) {
	if left >= right {
		return
	}
	partIndex := part(arr, left, right)
	quickSort(arr, left, partIndex-1)
	quickSort(arr, partIndex+1, right)
}

func part(arr []int, left, right int) int {
	// generate random index
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(right)
	var partIndex int
	if left > 0 {
		partIndex = random%left + left
	}

	// initialise pointer to signify last element of left partition
	lastOfLeft := left // at first, since left partiion is size 0, it is same
	// position as the left

	// move parition value to left index.
	arr[partIndex], arr[left] = arr[left], arr[partIndex]
	// start comparing values from index 1 to right and populate S1
	for i := left + 1; i <= right; i++ {
		if arr[i] < arr[left] {
			lastOfLeft++                                      //increase size
			arr[lastOfLeft], arr[i] = arr[i], arr[lastOfLeft] //populate s1
		}
	}
	// finally, shift partition value to the correct place by swapping it with the last element of left partition
	arr[lastOfLeft], arr[left] = arr[left], arr[lastOfLeft]

	return partIndex
}
