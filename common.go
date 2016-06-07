package main

import (
	"fmt"
)

func dumpExp(label string, bases [][]int, exp []int) {
	fmt.Printf("[%8s]:", label)

	for i, base := range bases {
		remain := base[1] - base[0]

		if exp[i+1] == -1 {
			fmt.Printf(" ----/%8d", remain)
		} else {
			fmt.Printf(" %4d/%d", exp[i+1], remain)
		}
	}

	fmt.Printf("\n")
}
