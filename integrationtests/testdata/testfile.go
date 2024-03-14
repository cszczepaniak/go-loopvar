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

func statementsInBetween() {
	for i, v := range []int{1, 2, 3} {
		index := i

		for range 10 {
		}

		go func() {}()
		val := v

		fmt.Println(index, val)
	}
}

func multiAssign() {
	for i, v := range []int{1, 2, 3} {
		idx, v := i, v

		fmt.Println(idx, v)
	}
}

func trickyMultiAssign() {
	for i, v := range []int{1, 2, 3} {
		_, anotherVar, v := i, 123, v

		fmt.Println(i, v, anotherVar)
	}
}

func wasABugBeforeGo122() {
	for i, v := range []int{1, 2, 3} {
		go func() {
			fmt.Println(i, v)
		}()

		idx := i
		val := v

		fmt.Println(idx, val)
	}
}

func muchWhitespace() {
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

func address() {
	vals := []*int{}
	for _, v := range []int{1, 2, 3} {
		v := v
		vals = append(vals, &v)
	}
	_ = vals
}

func addressRename() {
	vals := []*int{}
	for _, v := range []int{1, 2, 3} {
		myValue := v
		vals = append(vals, &myValue)
	}
	_ = vals
}

func renameInGoroutine() {
	for i, v := range []int{1, 2, 3} {
		index := i
		val := v

		go func() {
			fmt.Println(index, val)
		}()
	}
}

func kitchenSink() {
	for i1, v1 := range []int{1, 2, 3} {
		i1 := i1
		val1 := v1

		for i2, v2 := range []string{"foo", "bar"} {
			// I'm not sure how we should handle such captures... I don't expect they'll be common,
			// though.
			varA := i1
			anothaOne := v2

			go func() {
				for i3, v3 := range map[int]int{} {
					my3 := i3
					v3 := v3

					fmt.Println(
						i2,
						i1,
						val1,
						v1,
						my3,
						v3*(my3+i1-val1+v3),
						anothaOne,
						varA,
						v2,
					)
				}
			}()
		}
	}
}
