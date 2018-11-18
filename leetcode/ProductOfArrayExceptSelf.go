package leetcode

// https://leetcode.com/problems/product-of-array-except-self/

// O(n) with division
func productExceptSelf1(nums []int) []int {
	totalProduct := 1
	for _, val := range nums {
		totalProduct *= val
	}
	output := make([]int, len(nums))
	for i := range output {
		output[i] = totalProduct / nums[i]
	}
	return output
}

// O(n) by traversing twice in both directions and applying the correct multiplier.
//
// Left to Right:
// For a given index in nums array, we will first calculate the product of all
// values to the left of it and insert that product into the same index in output array
//
// Right to Left:
// Then, we will calculate the  product of all values to the right of this index.
// We then multiply this product with the value that exists in the output index,
// and replace that value with the result of this operation, giving us the final
// answer for each index.
func productExceptSelf2(nums []int) []int {
	output := make([]int, len(nums))

	leftMultiplier, rightMultiplier := 1, 1

	// left to right. initialize output values with the accumulated product, that is
	// the product of all values before this index.
	// If current index is the first index, then output value will just be by 1 first.
	for i := 0; i < len(nums); i++ {
		output[i] = leftMultiplier
		leftMultiplier *= nums[i]
	}

	// right to left. Correct the output values by multiplying each of them with the
	// accumulated product of all values to the right of the target index.
	// If current index is the last index, then the value will just be multiplied by 1.
	for j := len(nums) - 1; j >= 0; j-- {
		output[j] *= rightMultiplier
		rightMultiplier *= nums[j]
	}

	return output
}

// O(n^2)
func productExceptSelf3(nums []int) []int {
	productArray := make([]int, len(nums))
	for i := 0; i < len(nums); i++ {
		product := 1
		for j := 0; j < len(nums); j++ {
			if j == i {
				continue
			} else {
				product *= nums[j]
			}
		}
		productArray[i] = product
	}
	return productArray
}
