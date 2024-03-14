//go:build !exp

package testdata

import "fmt"

func slice() {
	// simple case
	for i, v := range []int{1, 2, 3} {
		i := i
		v := v

		fmt.Println(i, v)
	}

	// only thing in the loop
	for i, v := range []int{1, 2, 3} {
		i := i
		v := v
	}
}
