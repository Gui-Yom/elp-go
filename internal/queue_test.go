package internal

import (
	"fmt"
	"testing"
)

func TestSimple(t *testing.T) {
	queue := PriorityQueue{}
	queue.push("first", 1)
	queue.push("second", 2)
	queue.push("third", 3)
	fmt.Printf("%v\n", queue.pop())
	fmt.Printf("%v\n", queue.pop())
	fmt.Printf("%v\n", queue.pop())
	fmt.Printf("%v\n", queue.pop())
}

func TestPriority(t *testing.T) {
	queue := PriorityQueue{}
	queue.push("fourth", 4)
	queue.push("second", 2)
	queue.push("first", 1)
	queue.push("third", 3)
	queue.push("zeroth", 0)
	fmt.Printf("%v\n", queue.pop())
	fmt.Printf("%v\n", queue.pop())
	fmt.Printf("%v\n", queue.pop())
	fmt.Printf("%v\n", queue.pop())
	fmt.Printf("%v\n", queue.pop())
}
