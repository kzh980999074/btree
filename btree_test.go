package btree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewTree(t *testing.T) {
	root := buildBtreeRoot(4)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		InsertOne(root, i)
	}
	DeleteOne(root, 7)
	DeleteOne(root, 3)
	DeleteOne(root, 9)
	DeleteOne(root, 5)
	layers := [][]*Btree{}
	stack := []*Btree{root}
	for len(stack) != 0 {
		_stack := []*Btree{}
		layer := []*Btree{}
		for _, child := range stack {
			layer = append(layer, child)
			_stack = append(_stack, child.Childs...)
		}
		layers = append(layers, layer)
		stack = _stack
	}
	for _, vs := range layers {
		fmt.Println("--------------------------------")
		for _, v := range vs {
			fmt.Println(v)
		}
	}
}
