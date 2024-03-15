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

func trickyMultiAssign0() {
	for i, v := range []int{1, 2, 3} {
		_, anotherVar, val := i, 123, v

		fmt.Println(i, val, anotherVar)
	}
}

func trickyMultiAssign1() {
	for i, v := range []int{1, 2, 3} {
		_, index, anotherVar, val := 456, i, 123, v

		fmt.Println(index, val, anotherVar)
	}
}

func trickyMultiAssign2() {
	for _, v := range []int{1, 2, 3} {
		_, val := func() int { return 1 }(), v

		fmt.Println(val)
	}
}

func variableIsIncrementedLater() {
	for i, v := range []int{1, 2, 3} {
		incrementing, variable := i, v

		incrementing++

		fmt.Println(incrementing, variable)
	}
}

func variableIsDecrementedLater() {
	for i, v := range []int{1, 2, 3} {
		decrementing, variable := i, v

		decrementing--

		fmt.Println(decrementing, variable)
	}
}

func variableIsPlusAssigned() {
	for i, v := range []int{1, 2, 3} {
		incrementing, variable := i, v

		incrementing += 123

		fmt.Println(incrementing, variable)
	}
}

func variableIsMinusAssigned() {
	for i, v := range []int{1, 2, 3} {
		decrementing, variable := i, v

		decrementing -= 123

		fmt.Println(decrementing, variable)
	}
}

func variableIsMultiplyAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v

		index *= 123

		fmt.Println(index, variable)
	}
}

func variableIsDivideAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v

		index /= 123

		fmt.Println(index, variable)
	}
}

func variableIsModuloAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v

		index %= 123

		fmt.Println(index, variable)
	}
}

func variableIsAndAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v

		index &= 123

		fmt.Println(index, variable)
	}
}

func variableIsOrAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v

		index |= 123

		fmt.Println(index, variable)
	}
}

func variableIsXorAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v

		index ^= 123

		fmt.Println(index, variable)
	}
}

func variableIsShiftLeftAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v

		index <<= 123

		fmt.Println(index, variable)
	}
}

func variableIsShiftRightAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v

		index >>= 123

		fmt.Println(index, variable)
	}
}

func variableIsAndNotAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v

		index &^= 123

		fmt.Println(index, variable)
	}
}

func variableIsSimplyAssigned() {
	for i, v := range []int{1, 2, 3} {
		assigned, variable := i, v

		assigned = 123

		fmt.Println(assigned, variable)
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
