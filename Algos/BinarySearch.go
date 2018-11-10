package Algos

import (
	"fmt"
)

func main() {
	arr := []int{1, 2, 3, 4, 4, 5}
	if result := binarySearchIterative(arr, 4); result < 0 {
		fmt.Println("Target was not found in array!")
	} else {
		fmt.Println("Target was found in array index!", result)
	}
}

// binarySearchIterative performs binary search iteratively, assuming that
// it is sorted in ascending order.
func binarySearchIterative(arr []int, target int) int {
	left := 0
	right := len(arr) - 1
	mid := (left + right) / 2
	for left <= right {
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else if arr[mid] > target {
			right = mid - 1
		}
	}
	return -1 //not found
}

// binarySearchRecursive performs binary search recursively, assuming that
// it is sorted in ascending order.
func binarySearchRecursive(arr []int, left, right, target int) int {
	if left > right {
		return -1 //not found
	}
	mid := (left + right) / 2
	if arr[mid] == target {
		return mid
	}
	if arr[mid] > target {
		return binarySearchRecursive(arr, target, left, mid-1)
	} else {
		return binarySearchRecursive(arr, target, mid+1, right)
	}
}
