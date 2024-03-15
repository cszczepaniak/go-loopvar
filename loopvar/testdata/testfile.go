package testdata

import "fmt"

func slice() {
	for i, v := range []int{1, 2, 3} {
		i := i // want "found unnecessary loop variable capture"
		v := v // want "found unnecessary loop variable capture"

		fmt.Println(i, v)
	}
}

func statementsInBetween() {
	for i, v := range []int{1, 2, 3} {
		index := i // want "found unnecessary loop variable capture"

		for range 10 {
		}

		go func() {}()
		val := v // want "found unnecessary loop variable capture"

		fmt.Println(index, val)
	}
}

func multiAssign() {
	for i, v := range []int{1, 2, 3} {
		idx, v := i, v // want "found unnecessary loop variable capture"

		fmt.Println(idx, v)
	}
}

func trickyMultiAssign0() {
	for i, v := range []int{1, 2, 3} {
		_, anotherVar, val := i, 123, v // want "found unnecessary loop variable capture"

		fmt.Println(i, val, anotherVar)
	}
}

func trickyMultiAssign1() {
	for i, v := range []int{1, 2, 3} {
		_, index, anotherVar, val := 456, i, 123, v // want "found unnecessary loop variable capture"

		fmt.Println(index, val, anotherVar)
	}
}

func trickyMultiAssign2() {
	for _, v := range []int{1, 2, 3} {
		_, val := func() int { return 1 }(), v // want "found unnecessary loop variable capture"

		fmt.Println(val)
	}
}

func oneVariableIsIncrementedLaterOneIsNot() {
	for i, v := range []int{1, 2, 3} {
		incrementing := i
		variable := v // want "found unnecessary loop variable capture"

		incrementing++

		fmt.Println(incrementing, variable)
	}
}

func variableIsIncrementedLater() {
	for i, v := range []int{1, 2, 3} {
		incrementing, variable := i, v // want "found unnecessary loop variable capture"

		incrementing++

		fmt.Println(incrementing, variable)
	}
}

func variableIsDecrementedLater() {
	for i, v := range []int{1, 2, 3} {
		decrementing, variable := i, v // want "found unnecessary loop variable capture"

		decrementing--

		fmt.Println(decrementing, variable)
	}
}

func variableIsPlusAssigned() {
	for i, v := range []int{1, 2, 3} {
		incrementing, variable := i, v // want "found unnecessary loop variable capture"

		incrementing += 123

		fmt.Println(incrementing, variable)
	}
}

func variableIsMinusAssigned() {
	for i, v := range []int{1, 2, 3} {
		decrementing, variable := i, v // want "found unnecessary loop variable capture"

		decrementing -= 123

		fmt.Println(decrementing, variable)
	}
}

func variableIsMultiplyAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v // want "found unnecessary loop variable capture"

		index *= 123

		fmt.Println(index, variable)
	}
}

func variableIsDivideAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v // want "found unnecessary loop variable capture"

		index /= 123

		fmt.Println(index, variable)
	}
}

func variableIsModuloAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v // want "found unnecessary loop variable capture"

		index %= 123

		fmt.Println(index, variable)
	}
}

func variableIsAndAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v // want "found unnecessary loop variable capture"

		index &= 123

		fmt.Println(index, variable)
	}
}

func variableIsOrAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v // want "found unnecessary loop variable capture"

		index |= 123

		fmt.Println(index, variable)
	}
}

func variableIsXorAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v // want "found unnecessary loop variable capture"

		index ^= 123

		fmt.Println(index, variable)
	}
}

func variableIsShiftLeftAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v // want "found unnecessary loop variable capture"

		index <<= 123

		fmt.Println(index, variable)
	}
}

func variableIsShiftRightAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v // want "found unnecessary loop variable capture"

		index >>= 123

		fmt.Println(index, variable)
	}
}

func variableIsAndNotAssigned() {
	for i, v := range []int{1, 2, 3} {
		index, variable := i, v // want "found unnecessary loop variable capture"

		index &^= 123

		fmt.Println(index, variable)
	}
}

func variableIsSimplyAssigned() {
	for i, v := range []int{1, 2, 3} {
		assigned, variable := i, v // want "found unnecessary loop variable capture"

		assigned = 123

		fmt.Println(assigned, variable)
	}
}

func wasABugBeforeGo122() {
	for i, v := range []int{1, 2, 3} {
		go func() {
			fmt.Println(i, v)
		}()

		idx := i // want "found unnecessary loop variable capture"
		val := v // want "found unnecessary loop variable capture"

		fmt.Println(idx, val)
	}
}

func muchWhitespace() {
	for i, v := range []int{1, 2, 3} {

		i := i // want "found unnecessary loop variable capture"

		v := v // want "found unnecessary loop variable capture"

		fmt.Println(i, v)
	}
}

func aMap() {
	for k, v := range map[string]int{} {
		k := k // want "found unnecessary loop variable capture"
		v := v // want "found unnecessary loop variable capture"

		fmt.Println(k, v)
	}
}

func rename() {
	for i, v := range []int{1, 2, 3} {
		index := i // want "found unnecessary loop variable capture"
		val := v   // want "found unnecessary loop variable capture"

		fmt.Println(index, val)
	}
}

func address() {
	vals := []*int{}
	for _, v := range []int{1, 2, 3} {
		v := v // want "found unnecessary loop variable capture"
		vals = append(vals, &v)
	}
	_ = vals
}

func addressRename() {
	vals := []*int{}
	for _, v := range []int{1, 2, 3} {
		myValue := v // want "found unnecessary loop variable capture"
		vals = append(vals, &myValue)
	}
	_ = vals
}

func renameInGoroutine() {
	for i, v := range []int{1, 2, 3} {
		index := i // want "found unnecessary loop variable capture"
		val := v   // want "found unnecessary loop variable capture"

		go func() {
			fmt.Println(index, val)
		}()
	}
}

func kitchenSink() {
	for i1, v1 := range []int{1, 2, 3} {
		i1 := i1   // want "found unnecessary loop variable capture"
		val1 := v1 // want "found unnecessary loop variable capture"

		for i2, v2 := range []string{"foo", "bar"} {
			// I'm not sure how we should handle such captures... I don't expect they'll be common,
			// though.
			varA := i1
			anothaOne := v2 // want "found unnecessary loop variable capture"

			go func() {
				for i3, v3 := range map[int]int{} {
					my3 := i3 // want "found unnecessary loop variable capture"
					v3 := v3  // want "found unnecessary loop variable capture"

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
