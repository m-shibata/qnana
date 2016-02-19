package main

import (
	"encoding/json"
	"fmt"
)

type Stage3 struct {
	ApiFdam []float64 `json:"api_fdam"`
}

type Kouku struct {
	ApiStage3         Stage3 `json:"api_stage3"`
	ApiStage3Combined Stage3 `json:"api_stage3_combined"`
}

func (kouku Kouku) calcKoukuDamage(label string, hps1 []int, hps2 []int) {
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

type OpeningAtack struct {
	ApiFdam []float64 `json:"api_fdam"`
}

func (openingAtack OpeningAtack) calcOpeningAtackDamage(label string, hps []int) {
	fmt.Printf("[%8s]:", label)
	for i, v := range openingAtack.ApiFdam[1:] {
		hps[i] -= int(v)
		fmt.Printf(" %3d", int(v))
	}
	fmt.Printf("\n")
}

type Hougeki struct {
	ApiAtList []int         `json:"api_at_list"`
	ApiDamage []interface{} `json:"api_damage"`
	ApiDfList []interface{} `json:"api_df_list"`
}

func (hougeki Hougeki) calcHougekiDamage(label string, hps []int) {
	const leftMargin = "                |"
	fmt.Printf("[%8s]:\n", label)
	prevLeftSide := false
	for i, at := range hougeki.ApiAtList {
		if at == -1 {
			continue
		}
		margin := ""
		if at > len(hps) {
			if !prevLeftSide {
				margin = leftMargin
			} else {
				margin = " |"
			}
		} else if prevLeftSide {
			fmt.Printf("\n")
		}
		fmt.Printf("%s %2d", margin, at)
		var target []int
		switch t := hougeki.ApiDfList[i].(type) {
		case []interface{}:
			target = make([]int, len(t))
			for i, v := range t {
				target[i] = int(v.(float64))
			}
		default:
			target = nil
		}
		var damage []int
		switch t := hougeki.ApiDamage[i].(type) {
		case []interface{}:
			damage = make([]int, len(t))
			for i, v := range t {
				damage[i] = int(v.(float64))
			}
		default:
			damage = nil
		}
		fmt.Printf(" =>")
		for i, v := range target {
			if i > 0 {
				fmt.Printf("\n")
				if at > len(hps) {
					fmt.Printf("%s", leftMargin)
				}
				fmt.Printf("      ")
			}
			fmt.Printf(" %2d [%3d]", v, damage[i])
			if v > 0 && v <= len(hps) {
				hps[v-1] -= damage[i]
			}
		}
		if at > len(hps) {
			fmt.Printf("\n")
			prevLeftSide = false
		} else {
			prevLeftSide = true
		}
	}
	if prevLeftSide {
		fmt.Printf("\n")
	}
}

type Raigeki struct {
	ApiFdam []float64 `json:"api_fdam"`
}

func (raigeki Raigeki) calcRaigekiDamage(label string, hps []int) {
	fmt.Printf("[%8s]:", label)
	for i, v := range raigeki.ApiFdam[1:] {
		hps[i] -= int(v)
		fmt.Printf(" %3d", int(v))
	}
	fmt.Printf("\n")
}

type ApiReqCombinedBattleBattleWater struct {
	ApiShipKe         []int        `json:"api_ship_ke"`
	ApiKouku          Kouku        `json:"api_kouku"`
	ApiOpeningAtack   OpeningAtack `json:"api_opening_atack"`
	ApiHougeki1       Hougeki      `json:"api_hougeki1"`
	ApiHougeki2       Hougeki      `json:"api_hougeki2"`
	ApiHougeki3       Hougeki      `json:"api_hougeki3"`
	ApiRaigeki        Raigeki      `json:"api_raigeki"`
	ApiMaxhps         []int        `json:"api_maxhps"`
	ApiMaxhpsCombined []int        `json:"api_maxhps_combined"`
	ApiNowhps         []int        `json:"api_nowhps"`
	ApiNowhpsCombined []int        `json:"api_nowhps_combined"`
	ApiStageFlag      []int        `json:"api_stage_flag"`
	ApiOpeningFlag    int          `json:"api_opening_flag"`
	ApiHouraiFlag     []int        `json:"api_hourai_flag"`
}

type KcsapiApiReqCombinedBattleBattleWater struct {
	ApiData ApiReqCombinedBattleBattleWater `json:"api_data"`
	KcsapiBase
}

func dumpHps(label string, hps []int, maxhps []int) {
	fmt.Printf("%s", label)
	for i, hp := range hps {
		var flag string
		if hp > (maxhps[i+1] * 3 / 4) {
			flag = " "
		} else if hp > (maxhps[i+1] / 2) {
			flag = "-"
		} else if hp > (maxhps[i+1] / 4) {
			flag = "="
		} else {
			flag = "!"
		}
		fmt.Printf(" %3d/%3d%s", hp, maxhps[i+1], flag)
	}
	fmt.Printf("\n")
}

func handleApiReqCombinedBattleBattleWater(data []byte) error {
	var v KcsapiApiReqCombinedBattleBattleWater
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	enemy_size := len(v.ApiData.ApiShipKe) - 1

	deck1_hps := v.ApiData.ApiNowhps[1 : len(v.ApiData.ApiNowhps)-enemy_size]
	deck2_hps := v.ApiData.ApiNowhpsCombined[1:]

	if v.ApiData.ApiStageFlag[2] == 1 {
		v.ApiData.ApiKouku.calcKoukuDamage("Kouku", deck1_hps, deck2_hps)
	}
	if v.ApiData.ApiOpeningFlag == 1 {
		v.ApiData.ApiOpeningAtack.calcOpeningAtackDamage("Raigeki1", deck2_hps)
	}
	if v.ApiData.ApiHouraiFlag[0] == 1 {
		v.ApiData.ApiHougeki1.calcHougekiDamage("Hougeki1", deck1_hps)
	}
	if v.ApiData.ApiHouraiFlag[1] == 1 {
		v.ApiData.ApiHougeki2.calcHougekiDamage("Hougeki2", deck1_hps)
	}
	if v.ApiData.ApiHouraiFlag[2] == 1 {
		v.ApiData.ApiHougeki3.calcHougekiDamage("Hougeki3", deck2_hps)
	}
	if v.ApiData.ApiHouraiFlag[3] == 1 {
		v.ApiData.ApiRaigeki.calcRaigekiDamage("Raigeki2", deck2_hps)
	}

	dumpHps("Deck1", deck1_hps, v.ApiData.ApiMaxhps)
	dumpHps("Deck2", deck2_hps, v.ApiData.ApiMaxhpsCombined)

	return err
}

type ApiReqCombinedBattleLdAirbattle struct {
	ApiShipKe         []int `json:"api_ship_ke"`
	ApiKouku          Kouku `json:"api_kouku"`
	ApiMaxhps         []int `json:"api_maxhps"`
	ApiMaxhpsCombined []int `json:"api_maxhps_combined"`
	ApiNowhps         []int `json:"api_nowhps"`
	ApiNowhpsCombined []int `json:"api_nowhps_combined"`
	ApiStageFlag      []int `json:"api_stage_flag"`
}

type KcsapiApiReqCombinedBattleLdAirbattle struct {
	ApiData ApiReqCombinedBattleLdAirbattle `json:"api_data"`
	KcsapiBase
}

func handleApiReqCombinedBattleLdAirbattle(data []byte) error {
	var v KcsapiApiReqCombinedBattleLdAirbattle
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	enemy_size := len(v.ApiData.ApiShipKe) - 1

	deck1_hps := v.ApiData.ApiNowhps[1 : len(v.ApiData.ApiNowhps)-enemy_size]
	deck2_hps := v.ApiData.ApiNowhpsCombined[1:]

	if v.ApiData.ApiStageFlag[2] == 1 {
		v.ApiData.ApiKouku.calcKoukuDamage("Kouku", deck1_hps, deck2_hps)
	}

	dumpHps("Deck1", deck1_hps, v.ApiData.ApiMaxhps)
	dumpHps("Deck2", deck2_hps, v.ApiData.ApiMaxhpsCombined)

	return err
}
