package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Eventmap struct {
	ApiDmg      int `json:"api_dmg"`
	ApiMaxMaphp int `json:"api_max_maphp"`
	ApiNowMaphp int `json:"api_now_maphp"`
}

type EventId int

func (i EventId) String() string {
	switch i {
	case EventIdInit:
		return "Init"
	case EventIdUndef:
		return "Undef"
	case EventIdItem:
		return "Item"
	case EventIdWhirl:
		return "Whirl"
	case EventIdBattle:
		return "Battle"
	case EventIdBoss:
		return "Boss"
	case EventIdNone:
		return "None"
	case EventIdPlane:
		return "Plane"
	case EventIdGuard:
		return "Guard"
	case EventIdLand:
		return "Land"
	case EventIdLong:
		return "Long"
	}
	return strconv.Itoa(int(i))
}

const (
	EventIdInit EventId = iota
	EventIdUndef
	EventIdItem
	EventIdWhirl
	EventIdBattle
	EventIdBoss
	EventIdNone
	EventIdPlane
	EventIdGuard
	EventIdLand
	EventIdLong
)

type EventKind int

func (k EventKind) Label(i EventId) string {
	switch k {
	case EventKindNone:
		if i == EventIdPlane {
			return "Scout"
		}
		return "None"
	case EventKindNormal:
		if i == EventIdNone {
			return "None"
		}
		return "Normal"
	case EventKindNight:
		if i == EventIdNone {
			return "Fork"
		}
		return "Night"
	case EventKindDawn:
		return "Dawn"
	case EventKindAir:
		return "Air"
	case EventKindLong:
		return "Long"
	}
	return strconv.Itoa(int(i))
}

const (
	EventKindNone EventKind = iota
	EventKindNormal
	EventKindNight
	EventKindDawn
	EventKindAir
	EventKindLong
)

type MapNo int

const MapNos = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (i MapNo) String() string {

	if i < 1 || i > 26 {
		return "0"
	}

	return MapNos[i-1 : i]
}

type ApiReqMap struct {
	ApiAirsearch      interface{} `json:"api_airsearch"`
	ApiBosscellNo     MapNo       `json:"api_bosscell_no"`
	ApiBosscomp       int         `json:"api_bosscomp"`
	ApiColorNo        int         `json:"api_color_no"`
	ApiCommentKind    int         `json:"api_comment_kind"`
	ApiEventId        EventId     `json:"api_event_id"`
	ApiEventKind      EventKind   `json:"api_event_kind"`
	ApiEventmap       Eventmap    `json:"api_eventmap"`
	ApiFromNo         int         `json:"api_from_no"`
	ApiMapareaId      int         `json:"api_maparea_id"`
	ApiMapinfoNo      int         `json:"api_mapinfo_no"`
	ApiNext           int         `json:"api_next"`
	ApiNo             MapNo       `json:"api_no"`
	ApiProductionKind int         `json:"api_production_kind"`
	ApiRashinFlg      int         `json:"api_rashin_flg"`
	ApiRashinId       int         `json:"api_rashin_id"`
}

func (info ApiReqMap) dumpInfo() {
	fmt.Printf("[%2d - %2d]:", info.ApiMapareaId, info.ApiMapinfoNo)
	fmt.Printf(" Next: %s", info.ApiNo)
	fmt.Printf(" (%s/%s)", info.ApiEventId,
		info.ApiEventKind.Label(info.ApiEventId))
	if info.ApiNo == info.ApiBosscellNo {
		fmt.Printf(" BOSS")
	}
	if info.ApiEventmap.ApiMaxMaphp != 0 {
		fmt.Printf(" (HP: %4d/%4d)", info.ApiEventmap.ApiNowMaphp,
			info.ApiEventmap.ApiMaxMaphp)
	}

	fmt.Printf("\n")
}

type KcsapiApiReqMap struct {
	ApiData ApiReqMap `json:"api_data"`
	KcsapiBase
}

func handleApiReqMap(data []byte) error {
	var v KcsapiApiReqMap

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	v.ApiData.dumpInfo()
	return nil
}

func handleApiReqMapStart(data []byte) error {
	return handleApiReqMap(data)
}

func handleApiReqMapNext(data []byte) error {
	return handleApiReqMap(data)
}
