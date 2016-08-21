package main

import (
	"fmt"
	"strconv"
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
				if len(decksData[1]) < param.mvp2 {
					fmt.Printf(" / unknown")
				} else {
					fmt.Printf(" / %s", decksData[1][param.mvp2-1])
				}
			}
		}
	}
	fmt.Printf("\n")
}

func dumpExp(label string, bases [][]int, exp []int) {
	fmt.Printf("[%8s]:", label)

	for i, base := range bases {
		var remain int
		if len(base) < 2 {
			remain = 0
		} else {
			remain = base[1] - base[0]
		}

		if exp[i+1] == -1 {
			fmt.Printf(" ----/%d", remain)
		} else {
			fmt.Printf(" %4d/%d", exp[i+1], remain)
		}
	}

	fmt.Printf("\n")
}

type SkillLv int

func (i SkillLv) String() string {
	switch i {
	case 0:
		return ""
	case 1:
		return "|"
	case 2:
		return "||"
	case 3:
		return "|||"
	case 4:
		return "/"
	case 5:
		return "//"
	case 6:
		return "///"
	case 7:
		return ">>"
	}
	return strconv.Itoa(int(i))
}

type SlotItem struct {
	ApiId         int     `json:"api_id"`
	ApiLevel      int     `json:"api_level"`
	ApiLocked     int     `json:"api_locked"`
	ApiSlotitemId int     `json:"api_slotitem_id"`
	ApiAlv        SkillLv `json:"api_alv"`
}
