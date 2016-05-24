package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type MstShip struct {
	ApiId          int    `json:"api_id"`
	ApiAfterbull   int    `json:"api_afterbull"`
	ApiAfterfuel   int    `json:"api_afterfuel"`
	ApiAfterlv     int    `json:"api_afterlv"`
	ApiAftershipid string `json:"api_aftershipid"`
	ApiBacks       int    `json:"api_backs"`
	ApiBroken      []int  `json:"api_broken"`
	ApiBuidtime    int    `json:"api_buildtime"`
	ApiBullMax     int    `json:"api_bull_max"`
	ApiFuelMax     int    `json:"api_fuel_max"`
	ApiGetmes      string `json:"api_getmes"`
	ApiHoug        []int  `json:"api_houg"`
	ApiLeng        int    `json:"api_leng"`
	ApiLuck        []int  `json:"api_luck"`
	ApiMaxeq       []int  `json:"api_maxeq"`
	ApiName        string `json:"api_name"`
	ApiPowup       []int  `json:"api_powup"`
	ApiRaig        []int  `json:"api_raig"`
	ApiSlotNum     int    `json:"api_slot_num"`
	ApiSoku        int    `json:"api_soku"`
	ApiSortno      int    `json:"api_sortno"`
	ApiSouk        []int  `json:"api_souk"`
	ApiStype       int    `json:"api_stype"`
	ApiTaik        []int  `json:"api_taik"`
	ApiTyku        []int  `json:"api_tyku"`
	ApiVoicef      int    `json:"api_voicef"`
	ApiYomi        string `json:"api_yomi"`
}

type MstSlotitem struct {
	ApiCost     int    `json:"api_cost"`
	ApiDistance int    `json:"api_distance"`
	ApiAtap     int    `json:"api_atap"`
	ApiBakk     int    `json:"api_bakk"`
	ApiBaku     int    `json:"api_baku"`
	ApiBroken   []int  `json:"api_broken"`
	ApiHoug     int    `json:"api_houg"`
	ApiHouk     int    `json:"api_houk"`
	ApiHoum     int    `json:"api_houm"`
	ApiId       int    `json:"api_id"`
	ApiInfo     string `json:"api_info"`
	ApiLeng     int    `json:"api_leng"`
	ApiLuck     int    `json:"api_luck"`
	ApiName     string `json:"api_name"`
	ApiRaig     int    `json:"api_raig"`
	ApiRak      int    `json:"api_raik"`
	ApiRaim     int    `json:"api_raim"`
	ApiRare     int    `json:"api_rare"`
	ApiSakb     int    `json:"api_sakb"`
	ApiSaku     int    `json:"api_saku"`
	ApiSoku     int    `json:"api_soku"`
	ApiSortno   int    `json:"api_sortno"`
	ApiSouk     int    `json:"api_souk"`
	ApiTaik     int    `json:"api_taik"`
	ApiTais     int    `json:"api_tais"`
	ApiTyku     int    `json:"api_tyku"`
	ApiType     []int  `json:"api_type"`
	ApiUsebull  string `json:"api_usebull"`
}

type ShipData map[int]MstShip
type EquipData map[int]MstSlotitem

func (data ShipData) dumpShipNames(label string, ids []int) {
	fmt.Printf("[%8s]: ", label)

	list := make([]string, 0)
	for _, id := range ids {
		if id != -1 {
			list = append(list, data[id].ApiName)
		}
	}
	fmt.Println(strings.Join(list, " / "))
}

var shipData ShipData
var equipData EquipData

type ApiStart2 struct {
	ApiMstBgm               []interface{} `json:"api_mst_bgm"`
	ApiMstConst             interface{}   `json:"api_mst_const"`
	ApiMstEquipExslot       []int         `json:"api_mst_equip_exslot"`
	ApiMstFurniture         []interface{} `json:"api_mst_furniture"`
	ApiMstFurnituregraph    []interface{} `json:"api_mst_furnituregraph"`
	ApiMstItemShop          interface{}   `json:"api_mst_item_shop"`
	ApiMstMaparea           []interface{} `json:"api_mst_maparea"`
	ApiMstMapbgm            []interface{} `json:"api_mst_mapbgm"`
	ApiMstMapcell           []interface{} `json:"api_mst_mapcell"`
	ApiMstMapinfo           []interface{} `json:"api_mst_mapinfo"`
	ApiMstMission           []interface{} `json:"api_mst_mission"`
	ApiMstPayitem           []interface{} `json:"api_mst_payitem"`
	ApiMstShip              []MstShip     `json:"api_mst_ship"`
	ApiMstShipgraph         []interface{} `json:"api_mst_shipgraph"`
	ApiMstShipupgrade       []interface{} `json:"api_mst_shipupgrade"`
	ApiMstSlotitem          []MstSlotitem `json:"api_mst_slotitem"`
	ApiMstSlotitemEquiptype []interface{} `json:"api_mst_slotitem_equiptype"`
	ApiMstStype             []interface{} `json:"api_mst_stype"`
	ApiMstUseitem           []interface{} `json:"api_mst_useitem"`
	ApiRegisterStatus       int           `json:"api_register_status"`
}

func (data ApiStart2) updateShips() {
	var ship MstShip

	shipData = make(map[int]MstShip)

	for _, ship = range data.ApiMstShip {
		shipData[ship.ApiId] = ship
	}
}

func (data ApiStart2) updateEquips() {
	var equip MstSlotitem

	equipData = make(map[int]MstSlotitem)

	for _, equip = range data.ApiMstSlotitem {
		equipData[equip.ApiId] = equip
	}
}

func (data ApiStart2) loadLocal(uri string) error {
	f, err := ioutil.ReadFile(uri)
	if err != nil {
		return err
	}

	var v KcsapiApiStart2
	err = json.Unmarshal(f, &v)
	if err != nil {
		return err
	}

	v.ApiData.updateShips()
	v.ApiData.updateEquips()
	return err
}

type KcsapiApiStart2 struct {
	ApiData ApiStart2 `json:"api_data"`
	KcsapiBase
}

var dataSet ApiStart2

func handleApiStart2(data []byte) error {
	var v KcsapiApiStart2
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	v.ApiData.updateShips()
	v.ApiData.updateEquips()

	return err
}
