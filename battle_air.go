package main

import (
	"fmt"
	"strconv"
)

type SeikuId int

func (i SeikuId) String() string {
	switch i {
	case 0:
		return "均衡"
	case 1:
		return "確保"
	case 2:
		return "優勢"
	case 3:
		return "劣勢"
	case 4:
		return "喪失"
	}
	return strconv.Itoa(int(i))
}

type TouchPlane int

func (i TouchPlane) String() string {
	if i != -1 {
		if len(equipData) > int(i) {
			return equipData[int(i)].ApiName
		}
		return "OK"
	}
	return "NG"
}

type Stage1 struct {
	ApiDispSeiku  SeikuId      `json:"api_disp_seiku"`
	ApiECount     int          `json:"api_e_count"`
	ApiELostcount int          `json:"api_e_lostcount"`
	ApiFCount     int          `json:"api_f_count"`
	ApiFLostcount int          `json:"api_f_lostcount"`
	ApiTouchPlane []TouchPlane `json:"api_touch_plane"`
}

type Stage2 struct {
	ApiECount     int `json:"api_e_count"`
	ApiELostcount int `json:"api_e_lostcount"`
	ApiFCount     int `json:"api_f_count"`
	ApiFLostcount int `json:"api_f_lostcount"`
}

type Stage3 struct {
	ApiFdam []float64 `json:"api_fdam"`
}

type Kouku struct {
	ApiStage1         Stage1 `json:"api_stage1"`
	ApiStage2         Stage2 `json:"api_stage2"`
	ApiStage3         Stage3 `json:"api_stage3"`
	ApiStage3Combined Stage3 `json:"api_stage3_combined"`
}

func (kouku Kouku) calcKoukuDamage(label string, hps1 []int, hps2 []int) {
	fmt.Printf("[%7s0]: %-5s", label, kouku.ApiStage1.ApiDispSeiku)
	fmt.Printf("%10s / %10s / Touch\n", "All", "Bombers")
	fmt.Printf("            Friend %3d => %3d / %3d => %3d / %s\n",
		kouku.ApiStage1.ApiFCount, kouku.ApiStage1.ApiFCount-kouku.ApiStage1.ApiFLostcount,
		kouku.ApiStage2.ApiFCount, kouku.ApiStage2.ApiFCount-kouku.ApiStage2.ApiFLostcount,
		kouku.ApiStage1.ApiTouchPlane[0])
	fmt.Printf("            Enemy  %3d => %3d / %3d => %3d / %s\n",
		kouku.ApiStage1.ApiECount, kouku.ApiStage1.ApiECount-kouku.ApiStage1.ApiELostcount,
		kouku.ApiStage2.ApiECount, kouku.ApiStage2.ApiECount-kouku.ApiStage2.ApiELostcount,
		kouku.ApiStage1.ApiTouchPlane[1])
	fmt.Printf("[%7s1]:", label)
	for i, v := range kouku.ApiStage3.ApiFdam[1:] {
		hps1[i] -= int(v)
		fmt.Printf(" %3d", int(v))
	}
	fmt.Printf("\n")
	if hps2 == nil {
		return
	}
	fmt.Printf("[%7s2]:", label)
	for i, v := range kouku.ApiStage3Combined.ApiFdam[1:] {
		hps2[i] -= int(v)
		fmt.Printf(" %3d", int(v))
	}
	fmt.Printf("\n")
}
