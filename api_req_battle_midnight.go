package main

import (
	"encoding/json"
)

type ApiReqBattleMidnightBattle struct {
	ApiDeckId  int     `json:"api_deck_id"`
	ApiShipKe  []int   `json:"api_ship_ke"`
	ApiHougeki Hougeki `json:"api_hougeki"`
	ApiMaxhps  []int   `json:"api_maxhps"`
	ApiNowhps  []int   `json:"api_nowhps"`
}

type KcsapiApiReqBattleMidnightBattle struct {
	ApiData ApiReqBattleMidnightBattle `json:"api_data"`
	KcsapiBase
}

func handleApiReqBattleMidnightBattle(data []byte) error {
	var v KcsapiApiReqBattleMidnightBattle

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	currentDeckId = v.ApiData.ApiDeckId
	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, true)
	enemy_size := len(v.ApiData.ApiShipKe) - 1

	hps := v.ApiData.ApiNowhps[1 : len(v.ApiData.ApiNowhps)-enemy_size]

	v.ApiData.ApiHougeki.calcHougekiDamage("Hougeki", hps)

	dumpHps("Deck", hps, v.ApiData.ApiMaxhps)

	return err
}
