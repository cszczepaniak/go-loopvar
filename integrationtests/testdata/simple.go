package testdata

import "fmt"

func slice() {
	for i, v := range []int{1, 2, 3} {
		i := i
		v := v

		fmt.Println(i, v)
	}
}