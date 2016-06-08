package main

import (
	"encoding/json"
	"fmt"
)

type ApiReqSortieBattle struct {
	ApiDockId       int          `json:"api_dock_id"`
	ApiFormation    Formation    `json:"api_formation"`
	ApiShipKe       []int        `json:"api_ship_ke"`
	ApiKouku        Kouku        `json:"api_kouku"`
	ApiOpeningAtack OpeningAtack `json:"api_opening_atack"`
	ApiHougeki1     Hougeki      `json:"api_hougeki1"`
	ApiHougeki2     Hougeki      `json:"api_hougeki2"`
	ApiHougeki3     Hougeki      `json:"api_hougeki3"`
	ApiRaigeki      Raigeki      `json:"api_raigeki"`
	ApiMaxhps       []int        `json:"api_maxhps"`
	ApiNowhps       []int        `json:"api_nowhps"`
	ApiStageFlag    []int        `json:"api_stage_flag"`
	ApiOpeningFlag  int          `json:"api_opening_flag"`
	ApiHouraiFlag   []int        `json:"api_hourai_flag"`
}

type KcsapiApiReqSortieBattle struct {
	ApiData ApiReqSortieBattle `json:"api_data"`
	KcsapiBase
}

func handleApiReqSortieBattle(data []byte) error {
	var v KcsapiApiReqSortieBattle

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	currentDeckId = v.ApiData.ApiDockId
	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, true)
	enemy_size := len(v.ApiData.ApiShipKe) - 1

	hps := v.ApiData.ApiNowhps[1 : len(v.ApiData.ApiNowhps)-enemy_size]

	v.ApiData.ApiFormation.dumpFormation()
	if v.ApiData.ApiStageFlag[2] == 1 {
		v.ApiData.ApiKouku.calcKoukuDamage("Kouku", hps, nil)
	}
	if v.ApiData.ApiOpeningFlag == 1 {
		v.ApiData.ApiOpeningAtack.calcOpeningAtackDamage("Raigeki1", hps)
	}
	if v.ApiData.ApiHouraiFlag[0] == 1 {
		v.ApiData.ApiHougeki1.calcHougekiDamage("Hougeki1", hps)
	}
	if v.ApiData.ApiHouraiFlag[1] == 1 {
		v.ApiData.ApiHougeki2.calcHougekiDamage("Hougeki2", hps)
	}
	if v.ApiData.ApiHouraiFlag[2] == 1 {
		v.ApiData.ApiHougeki3.calcHougekiDamage("Hougeki3", hps)
	}
	if v.ApiData.ApiHouraiFlag[3] == 1 {
		v.ApiData.ApiRaigeki.calcRaigekiDamage("Raigeki2", hps)
	}

	dumpHps("Deck", hps, v.ApiData.ApiMaxhps)

	return err
}

type ApiReqSortieAirbattle struct {
	ApiDockId     int       `json:"api_dock_id"`
	ApiFormation  Formation `json:"api_formation"`
	ApiShipKe     []int     `json:"api_ship_ke"`
	ApiKouku      Kouku     `json:"api_kouku"`
	ApiKouku2     Kouku     `json:"api_kouku2"`
	ApiMaxhps     []int     `json:"api_maxhps"`
	ApiNowhps     []int     `json:"api_nowhps"`
	ApiStageFlag  []int     `json:"api_stage_flag"`
	ApiStageFlag2 []int     `json:"api_stage_flag2"`
}

type KcsapiApiReqSortieAirbattle struct {
	ApiData ApiReqSortieAirbattle `json:"api_data"`
	KcsapiBase
}

func handleApiReqSortieAirbattle(data []byte) error {
	var v KcsapiApiReqSortieAirbattle
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	currentDeckId = v.ApiData.ApiDockId
	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, true)
	enemy_size := len(v.ApiData.ApiShipKe) - 1

	hps := v.ApiData.ApiNowhps[1 : len(v.ApiData.ApiNowhps)-enemy_size]

	v.ApiData.ApiFormation.dumpFormation()
	if v.ApiData.ApiStageFlag[2] == 1 {
		v.ApiData.ApiKouku.calcKoukuDamage("KoukuA", hps, nil)
	}
	if v.ApiData.ApiStageFlag2[2] == 1 {
		v.ApiData.ApiKouku2.calcKoukuDamage("KoukuB", hps, nil)
	}

	dumpHps("Deck", hps, v.ApiData.ApiMaxhps)

	return err
}

type ApiReqSortieLdAirbattle struct {
	ApiDockId    int       `json:"api_dock_id"`
	ApiFormation Formation `json:"api_formation"`
	ApiShipKe    []int     `json:"api_ship_ke"`
	ApiKouku     Kouku     `json:"api_kouku"`
	ApiMaxhps    []int     `json:"api_maxhps"`
	ApiNowhps    []int     `json:"api_nowhps"`
	ApiStageFlag []int     `json:"api_stage_flag"`
}

type KcsapiApiReqSortieLdAirbattle struct {
	ApiData ApiReqSortieAirbattle `json:"api_data"`
	KcsapiBase
}

func handleApiReqSortieLdAirbattle(data []byte) error {
	var v KcsapiApiReqSortieLdAirbattle
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	currentDeckId = v.ApiData.ApiDockId
	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, true)
	enemy_size := len(v.ApiData.ApiShipKe) - 1

	hps := v.ApiData.ApiNowhps[1 : len(v.ApiData.ApiNowhps)-enemy_size]

	v.ApiData.ApiFormation.dumpFormation()
	if v.ApiData.ApiStageFlag[2] == 1 {
		v.ApiData.ApiKouku.calcKoukuDamage("Kouku", hps, nil)
	}

	dumpHps("Deck", hps, v.ApiData.ApiMaxhps)

	return err
}

type ApiReqSortieBattleresult struct {
	ApiGetExpLvup [][]int `json:"api_get_exp_lvup"`
	ApiGetShipExp []int   `json:"api_get_ship_exp"`
	ApiGetFlag    []int   `json:"api_get_flag"`
	ApiGetShip    GetShip `json:"api_get_ship"`
	ApiMvp        int     `json:"api_mvp"`
	ApiWinRank    string  `json:"api_win_rank"`
}

type KcsapiApiReqSortieBattleresult struct {
	ApiData ApiReqSortieBattleresult `json:"api_data"`
	KcsapiBase
}

func handleApiReqSortieBattleresult(data []byte) error {
	var v KcsapiApiReqSortieBattleresult

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	dumpRank(RankParams{label: "Rank", rank: v.ApiData.ApiWinRank,
		mvp1: v.ApiData.ApiMvp})
	dumpExp("Exp", v.ApiData.ApiGetExpLvup, v.ApiData.ApiGetShipExp)
	if v.ApiData.ApiGetFlag[1] == 1 {
		fmt.Printf("[Reunited]: %s %s\n", v.ApiData.ApiGetShip.ApiShipType, v.ApiData.ApiGetShip.ApiShipName)
	}

	return err
}
