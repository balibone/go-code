package hackerrank

// TreeNode ...
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// breadthFirstSearch is an iterative function to search for a leaf in a tree
// with a certain value. returning true if found
func breadthFirstSearch(head TreeNode, val int) bool {
	arr := []int{head.Val}
	// while array of values is not empty
	for len(arr) > 0 {
		// if val is found, return true
		if head.Val == val {
			return true
		}
		// if left exists, offer left
		if head.Left != nil {
			arr = append(arr, head.Left.Val)
		}
		// if right exists, offer right.
		if head.Right != nil {
			arr = append(arr, head.Right.Val)
		}
		// truncate head
		arr = arr[1:]
	}
}
