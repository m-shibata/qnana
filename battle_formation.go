package main

import (
	"fmt"
	"strconv"
)

type Formation []interface{}

func (f Formation) Formation(i int) string {
	if len(f) <= i {
		return "不明"
	}
	switch t := f[i].(type) {
	case string:
		switch t {
		case "1":
			return "単縦陣"
		case "2":
			return "複縦陣"
		case "3":
			return "輪形陣"
		case "4":
			return "梯形陣"
		case "5":
			return "単横陣"
		case "11":
			return "第1警戒航行序列"
		case "12":
			return "第2警戒航行序列"
		case "13":
			return "第3警戒航行序列"
		case "14":
			return "第4警戒航行序列"
		}
		return t
	case float64:
		switch t {
		case 1:
			return "単縦陣"
		case 2:
			return "複縦陣"
		case 3:
			return "輪形陣"
		case 4:
			return "梯形陣"
		case 5:
			return "単横陣"
		}
		return strconv.Itoa(int(t))
	}
	return "不明"
}

func (f Formation) Tactics() string {
	if len(f) < 3 {
		return "不明"
	}
	switch t := f[2].(type) {
	case float64:
		switch t {
		case 1:
			return "同航戦"
		case 2:
			return "反航戦"
		case 3:
			return "T字有利"
		case 4:
			return "T字不利"
		}
		return strconv.Itoa(int(t))
	}
	return "不明"
}

func (f Formation) dumpFormation() {
	fmt.Printf("[%8s]: %s vs. %s : %s\n", "Tactics",
		f.Formation(0), f.Formation(1), f.Tactics())
}
