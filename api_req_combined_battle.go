package main

import (
	"encoding/json"
	"fmt"
)

type Stage1 struct {
	ApiDispSeiku  int `json:"api_disp_seiku"`
	ApiECount     int `json:"api_e_count"`
	ApiELostcount int `json:"api_e_lostcount"`
	ApiFCount     int `json:"api_f_count"`
	ApiFLostcount int `json:"api_f_lostcount"`
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
	fmt.Printf("[%7s0]:", label)
	switch kouku.ApiStage1.ApiDispSeiku {
	case 0:
		fmt.Printf(" %-7s", "(B)")
	case 1:
		fmt.Printf(" %-7s", "(S)")
	case 2:
		fmt.Printf(" %-7s", "(A)")
	case 3:
		fmt.Printf(" %-7s", "(C)")
	case 4:
		fmt.Printf(" %-7s", "(D)")
	default:
		fmt.Printf(" %-7s", "(-)")
	}
	fmt.Printf("%10s / %10s\n", "All", "Bombers")
	fmt.Printf("            Friend %3d => %3d / %3d => %3d\n",
		kouku.ApiStage1.ApiFCount, kouku.ApiStage1.ApiFCount-kouku.ApiStage1.ApiFLostcount,
		kouku.ApiStage2.ApiFCount, kouku.ApiStage2.ApiFCount-kouku.ApiStage2.ApiFLostcount)
	fmt.Printf("            Enemy  %3d => %3d / %3d => %3d\n",
		kouku.ApiStage1.ApiECount, kouku.ApiStage1.ApiECount-kouku.ApiStage1.ApiELostcount,
		kouku.ApiStage2.ApiECount, kouku.ApiStage2.ApiECount-kouku.ApiStage2.ApiELostcount)
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
	fmt.Printf("[%8s]:", label)
	for i, hp := range hps {
		var flag string
		if maxhps[i+1] < 0 {
			fmt.Printf(" ---/--- ")
			continue
		} else if hp > (maxhps[i+1] * 3 / 4) {
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

type ApiReqCombinedBattleMidnightBattle struct {
	ApiShipKe         []int   `json:"api_ship_ke"`
	ApiHougeki        Hougeki `json:"api_hougeki"`
	ApiMaxhps         []int   `json:"api_maxhps"`
	ApiMaxhpsCombined []int   `json:"api_maxhps_combined"`
	ApiNowhps         []int   `json:"api_nowhps"`
	ApiNowhpsCombined []int   `json:"api_nowhps_combined"`
}

type KcsapiApiReqCombinedBattleMidnightBattle struct {
	ApiData ApiReqCombinedBattleMidnightBattle `json:"api_data"`
	KcsapiBase
}

func handleApiReqCombinedBattleMidnightBattle(data []byte) error {
	var v KcsapiApiReqCombinedBattleMidnightBattle
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	enemy_size := len(v.ApiData.ApiShipKe) - 1

	deck1_hps := v.ApiData.ApiNowhps[1 : len(v.ApiData.ApiNowhps)-enemy_size]
	deck2_hps := v.ApiData.ApiNowhpsCombined[1:]

	v.ApiData.ApiHougeki.calcHougekiDamage("Hougeki", deck2_hps)

	dumpHps("Deck1", deck1_hps, v.ApiData.ApiMaxhps)
	dumpHps("Deck2", deck2_hps, v.ApiData.ApiMaxhpsCombined)

	return err
}

type GetShip struct {
	ApiShipId   int    `json:"api_ship_id"`
	ApiShipName string `json:"api_ship_name"`
	ApiShipType string `json:"api_ship_type"`
}

type ApiReqCombinedBattleBattleresult struct {
	ApiGetFlag []int   `json:"api_get_flag"`
	ApiGetShip GetShip `json:"api_get_ship"`
	ApiWinRank string  `json:"api_win_rank"`
}

type KcsapiApiReqCombinedBattleBattleresult struct {
	ApiData ApiReqCombinedBattleBattleresult `json:"api_data"`
	KcsapiBase
}

func handleApiReqCombinedBattleBattleresult(data []byte) error {
	var v KcsapiApiReqCombinedBattleBattleresult
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	fmt.Printf("[%8s]: %s\n", "Rank", v.ApiData.ApiWinRank)
	if v.ApiData.ApiGetFlag[1] == 1 {
		fmt.Printf("[Reunited]: %s %s\n", v.ApiData.ApiGetShip.ApiShipType, v.ApiData.ApiGetShip.ApiShipName)
	}
	return err
}
