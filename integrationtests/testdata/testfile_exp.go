//go:build exp

package testdata

import "fmt"

func slice() {
	for i, v := range []int{1, 2, 3} {
		fmt.Println(i, v)
	}
}

func muchWhitespace() {
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

func kitchenSink() {
	for i1, v1 := range []int{1, 2, 3} {
		for i2, v2 := range []string{"foo", "bar"} {
			// I'm not sure how we should handle such captures... I don't expect they'll be common,
			// though.
			varA := i1
			go func() {
				for i3, v3 := range map[int]int{} {
					fmt.Println(
						i2,
						i1,
						v1,
						v1,
						i3,
						v3*(i3+i1-v1+v3),
						v2,
						varA,
						v2,
					)
				}
			}()
		}
	}
}
