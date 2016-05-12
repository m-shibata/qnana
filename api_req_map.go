package main

import (
	"encoding/json"
	"fmt"
)

type Eventmap struct {
	ApiDmg      int `json:"api_dmg"`
	ApiMaxMaphp int `json:"api_max_maphp"`
	ApiNowMaphp int `json:"api_now_maphp"`
}

type ApiReqMap struct {
	ApiAirsearch      interface{} `json:"api_airsearch"`
	ApiBosscellNo     int         `json:"api_bosscell_no"`
	ApiBosscomp       int         `json:"api_bosscomp"`
	ApiColorNo        int         `json:"api_color_no"`
	ApiCommentKind    int         `json:"api_comment_kind"`
	ApiEventId        int         `json:"api_event_id"`
	ApiEventKind      int         `json:"api_event_kind"`
	ApiEventmap       Eventmap    `json:"api_eventmap"`
	ApiFromNo         int         `json:"api_from_no"`
	ApiMapareaId      int         `json:"api_maparea_id"`
	ApiMapinfoNo      int         `json:"api_mapinfo_no"`
	ApiNext           int         `json:"api_next"`
	ApiNo             int         `json:"api_no"`
	ApiProductionKind int         `json:"api_production_kind"`
	ApiRashinFlg      int         `json:"api_rashin_flg"`
	ApiRashinId       int         `json:"api_rashin_id"`
}

func (info ApiReqMap) dumpInfo() {
	fmt.Printf("[%2d - %2d]:", info.ApiMapareaId, info.ApiMapinfoNo)
	fmt.Printf(" %2d => %2d", info.ApiFromNo, info.ApiNo)
	fmt.Printf(" (Event:%d Type:%d)", info.ApiEventId, info.ApiEventKind)
	if info.ApiNo == info.ApiBosscellNo {
		fmt.Printf(" BOSS map")
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
