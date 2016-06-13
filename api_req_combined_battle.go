package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type OpeningAtack struct {
	ApiFdam []float64 `json:"api_fdam"`
	ApiEdam []float64 `json:"api_edam"`
}

func (openingAtack OpeningAtack) calcOpeningAtackDamage(label string, dmg Damage, deck int) {
	fmt.Printf("[%8s]:", label)
	for i, v := range openingAtack.ApiFdam[1:] {
		dmg.deck[deck].dmg[i+1] += int(v)
		fmt.Printf(" %3d", int(v))
	}
	for i, v := range openingAtack.ApiEdam[1:] {
		dmg.enemy.dmg[i] += int(v)
	}
	fmt.Printf("\n")
}

type Hougeki struct {
	ApiAtList []int         `json:"api_at_list"`
	ApiDamage []interface{} `json:"api_damage"`
	ApiDfList []interface{} `json:"api_df_list"`
}

func (hougeki Hougeki) calcHougekiDamage(label string, dmg Damage, deck int) {
	const leftMargin = "                |"
	fmt.Printf("[%8s]:\n", label)
	prevLeftSide := false
	for i, at := range hougeki.ApiAtList {
		if at == -1 {
			continue
		}
		margin := ""
		if at >= len(dmg.deck[deck].dmg) {
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
				if at >= len(dmg.deck[deck].dmg) {
					fmt.Printf("%s", leftMargin)
				}
				fmt.Printf("      ")
			}
			fmt.Printf(" %2d [%3d]", v, damage[i])
			if v > 0 {
				if v < len(dmg.deck[deck].dmg) {
					dmg.deck[deck].dmg[v] += damage[i]
				} else {
					dmg.enemy.dmg[v-len(dmg.deck[deck].dmg)] += damage[i]
				}
			}
		}
		if at >= len(dmg.deck[deck].dmg) {
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
	ApiEdam []float64 `json:"api_edam"`
}

func (raigeki Raigeki) calcRaigekiDamage(label string, dmg Damage, deck int) {
	fmt.Printf("[%8s]:", label)
	for i, v := range raigeki.ApiFdam[1:] {
		dmg.deck[deck].dmg[i+1] += int(v)
		fmt.Printf(" %3d", int(v))
	}
	for i, v := range raigeki.ApiEdam[1:] {
		dmg.enemy.dmg[i] += int(v)
	}
	fmt.Printf("\n")
}

type ApiReqCombinedBattleBattle struct {
	ApiDeckId         string       `json:"api_deck_id"`
	ApiFormation      Formation    `json:"api_formation"`
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

type KcsapiApiReqCombinedBattleBattle struct {
	ApiData ApiReqCombinedBattleBattle `json:"api_data"`
	KcsapiBase
}

type ApiReqCombinedBattleBattleWater struct {
	ApiDeckId         string       `json:"api_deck_id"`
	ApiFormation      Formation    `json:"api_formation"`
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

func handleApiReqCombinedBattleBattle(data []byte) error {
	var v KcsapiApiReqCombinedBattleBattle
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	currentDeckId, _ = strconv.Atoi(v.ApiData.ApiDeckId)
	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, true)

	var damage Damage
	damage.init(v.ApiData.ApiNowhps, v.ApiData.ApiMaxhps, v.ApiData.ApiShipKe)
	damage.initCombined(v.ApiData.ApiNowhpsCombined, v.ApiData.ApiMaxhpsCombined)

	v.ApiData.ApiFormation.dumpFormation()
	if v.ApiData.ApiStageFlag[2] == 1 {
		v.ApiData.ApiKouku.calcKoukuDamage("Kouku", damage)
	}
	if v.ApiData.ApiOpeningFlag == 1 {
		v.ApiData.ApiOpeningAtack.calcOpeningAtackDamage("Raigeki1", damage, 1)
	}
	if v.ApiData.ApiHouraiFlag[0] == 1 {
		v.ApiData.ApiHougeki1.calcHougekiDamage("Hougeki1", damage, 1)
	}
	if v.ApiData.ApiHouraiFlag[1] == 1 {
		v.ApiData.ApiRaigeki.calcRaigekiDamage("Raigeki2", damage, 1)
	}
	if v.ApiData.ApiHouraiFlag[2] == 1 {
		v.ApiData.ApiHougeki3.calcHougekiDamage("Hougeki2", damage, 0)
	}
	if v.ApiData.ApiHouraiFlag[3] == 1 {
		v.ApiData.ApiHougeki2.calcHougekiDamage("Hougeki3", damage, 0)
	}

	damage.dumpHps()

	return err
}

func handleApiReqCombinedBattleBattleWater(data []byte) error {
	var v KcsapiApiReqCombinedBattleBattleWater
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	currentDeckId, _ = strconv.Atoi(v.ApiData.ApiDeckId)
	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, true)

	var damage Damage
	damage.init(v.ApiData.ApiNowhps, v.ApiData.ApiMaxhps, v.ApiData.ApiShipKe)
	damage.initCombined(v.ApiData.ApiNowhpsCombined, v.ApiData.ApiMaxhpsCombined)

	v.ApiData.ApiFormation.dumpFormation()
	if v.ApiData.ApiStageFlag[2] == 1 {
		v.ApiData.ApiKouku.calcKoukuDamage("Kouku", damage)
	}
	if v.ApiData.ApiOpeningFlag == 1 {
		v.ApiData.ApiOpeningAtack.calcOpeningAtackDamage("Raigeki1", damage, 1)
	}
	if v.ApiData.ApiHouraiFlag[0] == 1 {
		v.ApiData.ApiHougeki1.calcHougekiDamage("Hougeki1", damage, 0)
	}
	if v.ApiData.ApiHouraiFlag[1] == 1 {
		v.ApiData.ApiHougeki3.calcHougekiDamage("Hougeki2", damage, 0)
	}
	if v.ApiData.ApiHouraiFlag[2] == 1 {
		v.ApiData.ApiHougeki2.calcHougekiDamage("Hougeki3", damage, 1)
	}
	if v.ApiData.ApiHouraiFlag[3] == 1 {
		v.ApiData.ApiRaigeki.calcRaigekiDamage("Raigeki2", damage, 1)
	}

	damage.dumpHps()

	return err
}

type ApiReqCombinedBattleLdAirbattle struct {
	ApiDeckId         string    `json:"api_deck_id"`
	ApiFormation      Formation `json:"api_formation"`
	ApiShipKe         []int     `json:"api_ship_ke"`
	ApiKouku          Kouku     `json:"api_kouku"`
	ApiMaxhps         []int     `json:"api_maxhps"`
	ApiMaxhpsCombined []int     `json:"api_maxhps_combined"`
	ApiNowhps         []int     `json:"api_nowhps"`
	ApiNowhpsCombined []int     `json:"api_nowhps_combined"`
	ApiStageFlag      []int     `json:"api_stage_flag"`
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

	currentDeckId, _ = strconv.Atoi(v.ApiData.ApiDeckId)
	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, true)

	var damage Damage
	damage.init(v.ApiData.ApiNowhps, v.ApiData.ApiMaxhps, v.ApiData.ApiShipKe)
	damage.initCombined(v.ApiData.ApiNowhpsCombined, v.ApiData.ApiMaxhpsCombined)

	v.ApiData.ApiFormation.dumpFormation()
	if v.ApiData.ApiStageFlag[2] == 1 {
		v.ApiData.ApiKouku.calcKoukuDamage("Kouku", damage)
	}

	damage.dumpHps()

	return err
}

type ApiReqCombinedBattleMidnightBattle struct {
	ApiDeckId         string  `json:"api_deck_id"`
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

	currentDeckId, _ = strconv.Atoi(v.ApiData.ApiDeckId)
	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, true)

	var damage Damage
	damage.init(v.ApiData.ApiNowhps, v.ApiData.ApiMaxhps, v.ApiData.ApiShipKe)
	damage.initCombined(v.ApiData.ApiNowhpsCombined, v.ApiData.ApiMaxhpsCombined)

	v.ApiData.ApiHougeki.calcHougekiDamage("Hougeki", damage, 1)

	damage.dumpHps()

	return err
}

type GetShip struct {
	ApiShipId   int    `json:"api_ship_id"`
	ApiShipName string `json:"api_ship_name"`
	ApiShipType string `json:"api_ship_type"`
}

type ApiReqCombinedBattleBattleresult struct {
	ApiGetExpLvup         [][]int `json:"api_get_exp_lvup"`
	ApiGetExpLvupCombined [][]int `json:"api_get_exp_lvup_combined"`
	ApiGetShipExp         []int   `json:"api_get_ship_exp"`
	ApiGetShipExpCombined []int   `json:"api_get_ship_exp_combined"`
	ApiGetFlag            []int   `json:"api_get_flag"`
	ApiGetShip            GetShip `json:"api_get_ship"`
	ApiMvp                int     `json:"api_mvp"`
	ApiMvpCombined        int     `json:"api_mvp_combined"`
	ApiWinRank            string  `json:"api_win_rank"`
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

	dumpRank(RankParams{label: "Rank", rank: v.ApiData.ApiWinRank,
		mvp1: v.ApiData.ApiMvp, mvp2: v.ApiData.ApiMvpCombined})
	dumpExp("Exp1", v.ApiData.ApiGetExpLvup, v.ApiData.ApiGetShipExp)
	dumpExp("Exp2", v.ApiData.ApiGetExpLvupCombined,
		v.ApiData.ApiGetShipExpCombined)
	if v.ApiData.ApiGetFlag[1] == 1 {
		fmt.Printf("[Reunited]: %s %s\n", v.ApiData.ApiGetShip.ApiShipType, v.ApiData.ApiGetShip.ApiShipName)
	}
	return err
}
