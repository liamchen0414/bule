package mdd

import (
	"fmt"
	"testing"
)

func TestBasicMDD(test *testing.T) {
	fmt.Println("TestBasicMDDs")
	mdd, _ := ex2()
	mdd.PrintDOT()
}

// taken from Paper, Figure 2, Example 6
func TestXOR(test *testing.T) {
	mdd := Init()
	v7 := mdd.NewNode(1, []int{1, 0})
	v6 := mdd.NewNode(1, []int{0, 1})
	v5 := mdd.NewNode(2, []int{v6, v7})
	v4 := mdd.NewNode(2, []int{v7, v6})
	v3 := mdd.NewNode(2, []int{v4, v5})
	v2 := mdd.NewNode(2, []int{v5, v4})
	mdd.NewNode(4, []int{v2, v3})
	mdd.PrintDOT()
}

func ex1() (mdd MddStore, top int) {
	mdd = Init()
	v4 := mdd.NewNode(1, []int{0, 0})
	v5 := mdd.NewNode(1, []int{1, 1})
	v6 := mdd.NewNode(1, []int{0, 1})
	v2 := mdd.NewNode(2, []int{v4, v5})
	v3 := mdd.NewNode(2, []int{v5, v6})
	top = mdd.NewNode(3, []int{v2, v3})
	return
}

func ex2() (mdd MddStore, top int) {
	mdd = Init()
	v6 := mdd.NewNode(1, []int{1, 0})
	v7 := mdd.NewNode(1, []int{0, 1})
	v4 := mdd.NewNode(1, []int{v6, v7})
	v5 := mdd.NewNode(1, []int{v7, v6})
	v2 := mdd.NewNode(1, []int{v4, v5})
	v3 := mdd.NewNode(1, []int{v5, v4})
	top = mdd.NewNode(3, []int{v2, v3})
	return
}
