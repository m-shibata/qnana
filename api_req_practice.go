package main

import (
	"encoding/json"
)

type ApiReqPracticeBattle struct {
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

	enemy_size := len(v.ApiData.ApiShipKe) - 1

	hps := v.ApiData.ApiNowhps[1 : len(v.ApiData.ApiNowhps)-enemy_size]

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

func handleApiReqPracticeMidnightBattle(data []byte) error {
	var v KcsapiApiReqPracticeBattle

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	return err
}
