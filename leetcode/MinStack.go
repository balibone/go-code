// https://leetcode.com/problems/min-stack

package leetcode

// MinStack has a GetMin() function that runs in O(1) time. Made possible by using
// a 2nd stack to keep track of local minimums on every pop & push.
type MinStack struct {
	stack    []int
	minStack []int
}

/** initialize your data structure here. */
func Constructor() MinStack {
	return MinStack{
		stack:    make([]int, 0),
		minStack: make([]int, 0),
	}
}

func (this *MinStack) Push(x int) {
	// push to stack normally
	this.stack = append(this.stack, x)
	// push to min stack only if this new element is LTE the current top of min
	// stack, or the real stack is empty, hence this element will be the smallest anyway
	//
	// LTE is to allow pushing of duplicate entries the current global minimum.
	if (len(this.minStack) > 0 && x <= this.minStack[len(this.minStack)-1]) || len(this.minStack) == 0 {
		this.minStack = append(this.minStack, x)
	}
}

func (this *MinStack) Pop() {
	var actualTop int
	// pop normally.
	if len(this.stack) > 0 {
		actualTop = this.stack[len(this.stack)-1]
		this.stack = this.stack[:len(this.stack)-1]
		// if to be popped is also top of minStack, pop it too.
		// dont need to check if len(minStack) > 0, because it will definitely be.
		if this.minStack[len(this.minStack)-1] == actualTop {
			this.minStack = this.minStack[:len(this.minStack)-1]
		}
	}
}

func (this *MinStack) Top() int {
	// peek normally.
	if len(this.stack) > 0 {
		return this.stack[len(this.stack)-1]
	} else {
		return this.stack[0]
	}
}

func (this *MinStack) GetMin() int {
	// peek minStack normally.
	if len(this.minStack) > 0 {
		return this.minStack[len(this.minStack)-1]
	} else {
		return this.minStack[0]
	}
}
