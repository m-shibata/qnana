package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ApiReqKousyouCreateitem struct {
	ApiCreateFlag int      `json:"api_create_flag"`
	ApiMaterial   []int    `json:"api_material"`
	ApiShizaiFlag int      `json:"api_shizai_flag"`
	ApiSlotItem   SlotItem `json:"api_slot_item"`
	ApiType3      int      `json:"api_type3"`
	ApiUnsetslot  []int    `json:"api_unsetslot"`
}

type KcsapiApiReqKousyouCreateitem struct {
	ApiData ApiReqKousyouCreateitem `json:"api_data"`
	KcsapiBase
}

func handleApiReqKousyouCreateitem(data []byte) error {
	var v KcsapiApiReqKousyouCreateitem

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	item := "failed"
	if v.ApiData.ApiCreateFlag == 1 {
		if len(equipData) > v.ApiData.ApiSlotItem.ApiSlotitemId {
			item = equipData[v.ApiData.ApiSlotItem.ApiSlotitemId].ApiName
		} else {
			item = "unknown"
		}
	}
	fmt.Printf("[%8s]: %s\n", "Item", item)
	return nil
}

type ApiReqKousyouGetship struct {
	ApiShipId   int        `json:"api_ship_id"`
	ApiSlotItem []SlotItem `json:"api_slot_item"`
}

type KcsapiApiReqKousyouGetship struct {
	ApiData ApiReqKousyouGetship `json:"api_data"`
	KcsapiBase
}

func handleApiReqKousyouGetship(data []byte) error {
	var v KcsapiApiReqKousyouGetship

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	ship := "unknown"
	if len(shipData) > v.ApiData.ApiShipId {
		ship = shipData[v.ApiData.ApiShipId].ApiName
	}
	items := make([]string, 0)
	for _, item := range v.ApiData.ApiSlotItem {
		if len(equipData) > item.ApiSlotitemId {
			items = append(items, equipData[item.ApiSlotitemId].ApiName)
		}
	}
	fmt.Printf("[%8s]: %s", "Reunited", ship)
	if len(items) > 0 {
		fmt.Printf(" (%s)", strings.Join(items, ", "))
	}
	fmt.Printf("\n")
	return nil
}
