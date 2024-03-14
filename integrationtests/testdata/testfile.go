//go:build !exp

package testdata

import "fmt"

func slice() {
	for i, v := range []int{1, 2, 3} {
		i := i
		v := v

		fmt.Println(i, v)
	}
}

func aMap() {
	for k, v := range map[string]int{} {
		k := k
		v := v

		fmt.Println(k, v)
	}
}

func rename() {
	for i, v := range []int{1, 2, 3} {
		index := i
		val := v

		fmt.Println(index, val)
	}
}
