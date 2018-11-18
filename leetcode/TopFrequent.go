package leetcode

import (
	"fmt"
	"sort"
)

// https://leetcode.com/problems/top-k-frequent-elements

// ElementCount stores the data for each element
type ElementCount struct {
	element int
	count   int
}

// with sorting.
func topKFrequent1(nums []int, k int) []int {
	// create result arr and map.
	arr := make([]ElementCount, 0)
	intMap := make(map[int]int)
	// go through nums and map the count of each number O(N)
	for _, integer := range nums {
		intMap[integer]++
	}
	fmt.Println("intMap", intMap)
	// go through map and append to array.
	for integer, count := range intMap { //O(N)
		arr = append(arr, ElementCount{
			element: integer,
			count:   count,
		})
	}
	// sort the array in reverse order of element occurence
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].count > arr[j].count
	})
	// return the first k elements in arr.
	truncated := arr[:k]
	answer := make([]int, k)
	for i := 0; i < k; i++ {
		answer[i] = truncated[i].element
	}
	return answer
}

// without sorting
func topKFrequent2(nums []int, k int) []int {
	// create result arr and map.
	arr := make([]ElementCount, 0)
	intMap := make(map[int]int)
	// go through nums and map the count of each number O(N)
	for _, integer := range nums {
		intMap[integer]++
	}
	fmt.Println("intMap", intMap)
	// go through map and extract out each number.
	for integer, count := range intMap { //O(N)
		this := ElementCount{
			element: integer,
			count:   count,
		}
		// if arr is empty at first, then just append and continue
		if len(arr) == 0 {
			arr = append(arr, this)
			fmt.Println("first", arr)
			continue
		}
		// O(n)
		arr = insertElementCount(this, arr)
		fmt.Println("after insert", arr)
	}
	// return the first k elements in arr.
	truncated := arr[:k]
	answer := make([]int, k)
	for i := 0; i < k; i++ {
		answer[i] = truncated[i].element
	}
	return answer
}

// insert element after the first index starting from the left where its count is LTE
func insertElementCount(elem ElementCount, arr []ElementCount) []ElementCount {
	fmt.Println("arr", arr)
	result := make([]ElementCount, 0)
	var inserted bool
	for i, el := range arr {
		if elem.count >= el.count {
			result = append(arr[:i], append([]ElementCount{elem}, arr[i:]...)...)
			inserted = true
			break
		}
	}
	// this guy is the smallest count so far. so append
	if inserted == false {
		result = append(arr, elem)
	}
	fmt.Println("result", result)
	return result
}
