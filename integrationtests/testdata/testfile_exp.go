//go:build exp

package testdata

import "fmt"

func slice() {
	for i, v := range []int{1, 2, 3} {
		fmt.Println(i, v)
	}
}

func aMap() {
	for k, v := range map[string]int{} {
		fmt.Println(k, v)
	}
}

func rename() {
	for i, v := range []int{1, 2, 3} {
		fmt.Println(i, v)
	}
}

func address() {
	vals := []*int{}
	for _, v := range []int{1, 2, 3} {
		vals = append(vals, &v)
	}
	_ = vals
}

func addressRename() {
	vals := []*int{}
	for _, v := range []int{1, 2, 3} {
		vals = append(vals, &v)
	}
	_ = vals
}

func renameInGoroutine() {
	for i, v := range []int{1, 2, 3} {
		go func() {
			fmt.Println(i, v)
		}()
	}
}
