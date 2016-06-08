package main

import (
	"encoding/json"
	"fmt"
)

type SlotItem struct {
	ApiId         int `json:"api_id"`
	ApiSlotitemId int `json:"api_slotitem_id"`
}

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
