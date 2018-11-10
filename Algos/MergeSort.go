package Algos

import "fmt"

func main() {
	arr := []int{3, 2, 5, 4}
	fmt.Println("Before:", arr)
	mergeSort(arr, 0, len(arr)-1)
	fmt.Println("After:", arr)
}

// mergeSort is a recursive function that can be broken down into 3 main behaviours:
//
// Divide:
// 1st, it will recursively divide the array into "left halves" until only 1 element remains.
// Thereafter, it will attempt mergeSort on the right half of the last sub array which
// had more than 1 element.
//
// Conquer:
// After finding out that this right half also only has 1 element left, a call to merge is made,
// which will merge and sort the 2 halves using the help of a temporary array.
//
// Looping back up the recursion stack:
// After the first conquer completes, the next divide will be carried out to next smallest right half.
// This right half goes through the same process of divide and conquer until finally,
// It is time for the result of the very first conquer and this sorted right half to be conquered.
// This continues for the rest of the array.
func mergeSort(arr []int, left, right int) {
	if left == right {
		// When this is true, means only 1 element left in this subarray,
		// so don't bother dividing anymore
		// It is impossible for left to overshoot right (i.e. left > right) because
		// (left+right)/2 will never be > right.
		// Maximum value that left can reach is the value of right.
		return
	}
	mid := (left + right) / 2
	mergeSort(arr, left, mid)
	mergeSort(arr, mid+1, right)
	merge(arr, left, mid, right)
}

func merge(arr []int, left, mid, right int) {
	// start of left half
	leftStart := left
	// start of right half
	rightStart := mid + 1
	// helper array to store values in sorted order
	tempArray := make([]int, right-left+1)
	tempPointer := 0
	// start merging
	for leftStart <= mid && rightStart <= right {
		if arr[leftStart] <= arr[rightStart] {
			tempArray[tempPointer] = arr[leftStart]
			leftStart++
		} else {
			tempArray[tempPointer] = arr[rightStart]
			rightStart++
		}
		tempPointer++
	}
	// if left half wasn't exhausted, exhaust it.
	for leftStart <= mid {
		tempArray[tempPointer] = arr[leftStart]
		tempPointer++
		leftStart++
	}
	// if right half wasn't exhausted, exhaust it.
	for rightStart <= right {
		tempArray[tempPointer] = arr[rightStart]
		tempPointer++
		rightStart++
	}
	// only 1 of the previous 2 for loops will run at any given time.
	// paste elements back into original array, but now in sorted order from tempArray
	for i, v := range tempArray {
		arr[left+i] = v
	}
}
