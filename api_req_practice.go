package main

import (
	"encoding/json"
)

type ApiReqPracticeBattle struct {
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

type KcsapiApiReqPracticeBattle struct {
	ApiData ApiReqPracticeBattle `json:"api_data"`
	KcsapiBase
}

func handleApiReqPracticeBattle(data []byte) error {
	var v KcsapiApiReqPracticeBattle

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	currentDeckId = v.ApiData.ApiDockId
	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, false)

	var damage Damage
	damage.init(v.ApiData.ApiNowhps, v.ApiData.ApiMaxhps, v.ApiData.ApiShipKe)

	v.ApiData.ApiFormation.dumpFormation()
	if v.ApiData.ApiStageFlag[2] == 1 {
		v.ApiData.ApiKouku.calcKoukuDamage("Kouku", damage)
	}
	if v.ApiData.ApiOpeningFlag == 1 {
		v.ApiData.ApiOpeningAtack.calcOpeningAtackDamage("Raigeki1", damage, 0)
	}
	if v.ApiData.ApiHouraiFlag[0] == 1 {
		v.ApiData.ApiHougeki1.calcHougekiDamage("Hougeki1", damage, 0)
	}
	if v.ApiData.ApiHouraiFlag[1] == 1 {
		v.ApiData.ApiHougeki2.calcHougekiDamage("Hougeki2", damage, 0)
	}
	if v.ApiData.ApiHouraiFlag[2] == 1 {
		v.ApiData.ApiHougeki3.calcHougekiDamage("Hougeki3", damage, 0)
	}
	if v.ApiData.ApiHouraiFlag[3] == 1 {
		v.ApiData.ApiRaigeki.calcRaigekiDamage("Raigeki2", damage, 0)
	}

	damage.dumpHps()

	return err
}

type ApiReqPracticeMidnightBattle struct {
	ApiShipKe  []int   `json:"api_ship_ke"`
	ApiHougeki Hougeki `json:"api_hougeki"`
	ApiMaxhps  []int   `json:"api_maxhps"`
	ApiNowhps  []int   `json:"api_nowhps"`
}

type KcsapiApiReqPracticeMidnightBattle struct {
	ApiData ApiReqPracticeMidnightBattle `json:"api_data"`
	KcsapiBase
}

func handleApiReqPracticeMidnightBattle(data []byte) error {
	var v KcsapiApiReqPracticeMidnightBattle

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, false)

	var damage Damage
	damage.init(v.ApiData.ApiNowhps, v.ApiData.ApiMaxhps, v.ApiData.ApiShipKe)
	v.ApiData.ApiHougeki.calcHougekiDamage("Hougeki", damage, 0)
	damage.dumpHps()

	return err
}

type ApiReqPracticeBattleResult struct {
	ApiGetExpLvup [][]int `json:"api_get_exp_lvup"`
	ApiGetShipExp []int   `json:"api_get_ship_exp"`
	ApiMvp        int     `json:"api_mvp"`
	ApiWinRank    string  `json:"api_win_rank"`
}

type KcsapiApiReqPracticeBattleResult struct {
	ApiData ApiReqPracticeBattleResult `json:"api_data"`
	KcsapiBase
}

func handleApiReqPracticeBattleResult(data []byte) error {
	var v KcsapiApiReqPracticeBattleResult

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	dumpRank(RankParams{label: "Rank", rank: v.ApiData.ApiWinRank,
		mvp1: v.ApiData.ApiMvp})
	dumpExp("Exp", v.ApiData.ApiGetExpLvup, v.ApiData.ApiGetShipExp)

	return err
}
