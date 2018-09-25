package Algos

func main() {
	arr := []int{1, 2, 3, 55132, 2313, 4}
	mergeSort(arr, 0, len(arr)-1)
}

func mergeSort(arr []int, left, right int) {
	if left >= right {
		return
	}
	mid := (left + right) / 2
	mergeSort(arr, left, mid)
	mergeSort(arr, mid+1, right)
	merge(arr, left, right)
}

func merge(arr []int, left, right int) {
	leftStart := left
	mid := (left + right) / 2
	rightStart := mid + 1
	helper := make([]int, right-left+1)
	for leftStart <= mid && rightStart <= right {
		if arr[leftStart] >= arr[rightStart] {
			helper = append(helper, arr[rightStart])
			rightStart++
		} else {
			helper = append(helper, arr[leftStart])
			leftStart++
		}
	}
	for leftStart <= mid {
		helper = append(helper, arr[leftStart])
		leftStart++
	}
	for rightStart <= right {
		helper = append(helper, arr[rightStart])
		rightStart++
	}
	copy(arr[left:right+1], helper)
}
