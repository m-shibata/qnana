package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

var itemsData map[int]SlotItem

type KcsapiApiGetMemberSlotItem struct {
	ApiData []SlotItem `json:"api_data"`
	KcsapiBase
}

func handleApiGetMemberSlotItem(data []byte) error {
	var v KcsapiApiGetMemberSlotItem
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	itemsData = make(map[int]SlotItem, len(v.ApiData))

    for _, item := range v.ApiData {
		itemsData[item.ApiId] = item
	}

	deck_alv := make([][]string, len(decksData))
	for i, deck := range decksData {
		var alv_list []string
		for j, ship := range deck {
			var alv []string
			for _, slot := range ship.ApiSlot {
				if item, ok := itemsData[slot]; ok {
					if item.ApiAlv != 0 {
						if equip, ok := equipData[item.ApiSlotitemId]; ok {
							alv = append(alv, fmt.Sprintf("%s[%3s]", equip.ApiName, item.ApiAlv))
						}
					}
				}
			}
			if len(alv) > 0 {
				alv_list = append(alv_list, fmt.Sprintf("%d: %s", j, strings.Join(alv, " / ")))
			}
		}
		if len(alv_list) > 0 {
			deck_alv[i] = alv_list
		}
	}

	for i, deck := range deck_alv {
		if len(deck) > 0 {
			fmt.Printf("Deck%d:\n", i)
			for _, alv := range deck {
				if len(alv) > 0 {
					fmt.Printf("  %s\n", alv)
				}
			}
		}
	}
	return err
}

