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
	label := param.label

	if param.mvp1 > 0 {
		label += "/MVP"
	}
	fmt.Printf("[%8s]: %s", label, param.rank)

	if param.mvp1 > 0 {
		if len(decksData) < currentDeckId {
			fmt.Printf(" / unknown deck (%d)", currentDeckId)
		} else if len(decksData[currentDeckId-1]) < param.mvp1 {
			fmt.Printf(" / unknown")
		} else {
			fmt.Printf(" / %s", decksData[currentDeckId-1][param.mvp1-1])
			if param.mvp2 > 0 {
				if len(decksData[currentDeckId-1]) < param.mvp2 {
					fmt.Printf(", unknown")
					fmt.Printf(", %s", decksData[currentDeckId-1][param.mvp2-1])
				}
			}
		}
	}
	fmt.Printf("\n")
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
