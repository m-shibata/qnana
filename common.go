package main

import (
	"fmt"
)

type RankParams struct {
	label      string
	rank       string
	mvp1, mvp2 int
}

func dumpRank(param RankParams) {

	if param.mvp1 <= 0 {
		fmt.Printf("[%8s]: %s\n", param.label, param.rank)
	} else if param.mvp2 <= 0 {
		fmt.Printf("[%8s]: %s(%d)\n", param.label+"/MVP",
			param.rank, param.mvp1)
	} else {
		fmt.Printf("[%8s]: %s(%d, %d)\n", param.label+"/MVP",
			param.rank, param.mvp1, param.mvp2)
	}
}

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
