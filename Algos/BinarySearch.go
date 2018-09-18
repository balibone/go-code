package main

import (
	"fmt"
)

func main() {
	arr := []int{1, 2, 3, 4, 4, 5}
	if result := binarySearch(arr, 4); result < 0 {
		fmt.Println("Target was not found in array!")
	} else {
		fmt.Println("Target was found in array index!", result)
	}
}

func binarySearch(arr []int, target int) int {
	return binarySearcher(arr, target, 0, len(arr)-1)
}

// this binary search assumes that the given array is sorted in ascending order.
func binarySearcher(arr []int, target, left, right int) int {
	if left > right {
		return -1 //overshot so not found
	}
	mid := (right + left) / 2
	if arr[mid] == target { //found
		return mid
	}
	if arr[mid] > target {
		return binarySearcher(arr, target, left, mid-1)
	} else {
		return binarySearcher(arr, target, mid+1, right)
	}
}
