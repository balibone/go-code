package Algos

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	arr := []int{1, 3, 0, 2, 5, 19, -1}
	fmt.Println("Before", arr)
	quickSort(arr, 0, len(arr)-1)
	fmt.Println("After:", arr)
}

// once quickSort is run, every element will have all smaller elements to the left and all larger elements to the right (aka sorted).
func quickSort(arr []int, left, right int) {
	if left >= right {
		return
	}
	pivotIndex := partitionArray(arr, left, right)
	quickSort(arr, left, pivotIndex-1)
	quickSort(arr, pivotIndex+1, right)
}

func partitionArray(arr []int, left, right int) int {
	// use a random index as partition index
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(right)
	var partIndex int
	if left > 0 {
		partIndex = random%left + left
	}

	// initialise pointer to signify last element of left partition
	// at first, since left partiion is size 0, it is same
	// position as the left
	lastOfLeft := left

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
