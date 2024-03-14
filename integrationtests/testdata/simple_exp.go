package testdata

import "fmt"

func slice() {
	for i, v := range []int{1, 2, 3} {
		fmt.Println(i, v)
	}
}
