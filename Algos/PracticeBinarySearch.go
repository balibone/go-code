package Algos

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

func binarySearchRecursive(arr []int, left, right, target int) int {
	if left > right {
		return -1 //not found
	}
	mid := (left + right) / 2
	if arr[mid] == target {
		return mid
	}
	if arr[mid] > target {
		return binarySearcher(arr, target, left, mid-1)
	} else {
		return binarySearcher(arr, target, mid+1, right)
	}
}
