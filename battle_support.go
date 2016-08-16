package main

import (
	"fmt"
)

type SupportAiratack struct {
	ApiPlaneFrom      [][]int `json:"api_plane_from"`
	ApiDeckId         int     `json:"api_deck_id"`
	ApiShipId         []int   `json:"api_ship_id"`
	ApiStage1         Stage1  `json:"api_stage1"`
	ApiStage2         Stage2  `json:"api_stage2"`
	ApiStage3         Stage3  `json:"api_stage3"`
	ApiStageFlag      []int   `json:"api_stage_flag"`
	ApiUndressingFlag []int   `json:"api_undressing_flag"`
}

type SupportHourai struct {
	ApiClList         []int     `json:"api_cl_list"`
	ApiDamage         []float64 `json:"api_damage"`
	ApiDeckId         int       `json:"api_deck_id"`
	ApiShipId         []int     `json:"api_ship_id"`
	ApiUndressingFlag []int     `json:"api_undressing_flag"`
}

type SupportInfo struct {
	ApiSupportAiratack SupportAiratack `json:"api_support_airatack"`
	ApiSupportHourai   SupportHourai   `json:"api_support_hourai"`
}

func (support SupportInfo) calcSupportDamage(flag int, dmg Damage) {
	if flag == 1 {
		fmt.Printf("[%8s]: 空撃\n", "Support")
		if support.ApiSupportAiratack.ApiStageFlag[2] == 1 {
			fmt.Printf("[%7s0]:", "Kouku")
			fmt.Printf("%10s / %10s\n", "All", "Bombers")
			fmt.Printf("      Friend %3d => %3d / %3d => %3d\n",
				support.ApiSupportAiratack.ApiStage1.ApiFCount,
				support.ApiSupportAiratack.ApiStage1.ApiFCount-support.ApiSupportAiratack.ApiStage1.ApiFLostcount,
				support.ApiSupportAiratack.ApiStage2.ApiFCount,
				support.ApiSupportAiratack.ApiStage2.ApiFCount-support.ApiSupportAiratack.ApiStage2.ApiFLostcount)
			fmt.Printf("[%7s1]: (Enemy)", "Kouku")
			for i, v := range support.ApiSupportAiratack.ApiStage3.ApiEdam[1:] {
				dmg.enemy.dmg[i] += int(v)
				fmt.Printf(" %3d", int(v))
			}
			fmt.Printf("\n")
		}
	} else if flag == 2 || flag == 3 {
		fmt.Printf("[%8s]: 砲雷撃 ", "Support")
		for i, v := range support.ApiSupportHourai.ApiDamage[1:] {
			dmg.enemy.dmg[i] += int(v)
			fmt.Printf(" %3d", int(v))
		}
		fmt.Printf("\n")
	}
}
