package main

import (
	"fmt"
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
	pivot := arr[left]
	leftHalfLast := left
	for i := left + 1; i <= right; i++ {
		if arr[i] < pivot { // place it in S1. make sure to swap the old first element of "S2" (aka the new last element of S1) with this element
			leftHalfLast++
			swap(arr, leftHalfLast, i)
		}
		// else, let everything stay there, because they are already in S2
	}
	// swap pivot with element in leftHalfLast so that pivot is now in the right position
	swap(arr, leftHalfLast, left)
	return leftHalfLast
}

func swap(arr []int, idxA, idxB int) {
	arr[idxA], arr[idxB] = arr[idxB], arr[idxA]
}
