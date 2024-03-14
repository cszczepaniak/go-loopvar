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
