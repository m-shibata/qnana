package main

import (
	"encoding/json"
	"strconv"
)

type ApiReqBattleMidnightBattle struct {
	ApiDeckId  string  `json:"api_deck_id"`
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

	currentDeckId, _ = strconv.Atoi(v.ApiData.ApiDeckId)
	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, true)

	var damage Damage
	damage.init(v.ApiData.ApiNowhps, v.ApiData.ApiMaxhps, v.ApiData.ApiShipKe)
	v.ApiData.ApiHougeki.calcHougekiDamage("Hougeki", damage, 0)
	damage.dumpHps()

	return err
}

type ApiReqBattleMidnightSpMidnight struct {
	ApiDeckId    int       `json:"api_deck_id"`
	ApiFormation Formation `json:"api_formation"`
	ApiShipKe    []int     `json:"api_ship_ke"`
	ApiHougeki   Hougeki   `json:"api_hougeki"`
	ApiMaxhps    []int     `json:"api_maxhps"`
	ApiNowhps    []int     `json:"api_nowhps"`
}

type KcsapiApiReqBattleMidnightSpMidnight struct {
	ApiData ApiReqBattleMidnightSpMidnight `json:"api_data"`
	KcsapiBase
}

func handleApiReqBattleMidnightSpMidnight(data []byte) error {
	var v KcsapiApiReqBattleMidnightSpMidnight

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	currentDeckId = v.ApiData.ApiDeckId
	shipData.dumpShipNames("Enemy", v.ApiData.ApiShipKe, true)

	var damage Damage
	damage.init(v.ApiData.ApiNowhps, v.ApiData.ApiMaxhps, v.ApiData.ApiShipKe)
	v.ApiData.ApiHougeki.calcHougekiDamage("Hougeki", damage, 0)
	damage.dumpHps()

	return err
}
