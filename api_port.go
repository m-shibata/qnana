package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

type Basic struct {
	ApiLevel int `json:"api_level"`
}

type Ship struct {
	ApiId       int     `json:"api_id"`
	ApiShipId   int     `json:"api_ship_id"`
	ApiCond     float64 `json:"api_cond"`
	ApiKaryoku  []int   `json:"api_karyoku"`
	ApiRaisou   []int   `json:"api_raisou"`
	ApiTaiku    []int   `json:"api_taiku"`
	ApiSoukou   []int   `json:"api_soukou"`
	ApiLucky    []int   `json:"api_lucky"`
	ApiKaihi    []int   `json:"api_kaihi"`
	ApiSakuteki []int   `json:"api_sakuteki"`
	ApiTaisen   []int   `json:"api_taisen"`
	ApiKyouka   []int   `json:"api_kyouka"`
	ApiSlot     []int   `json:"api_slot"`
	ApiSlotnum  int     `json:"api_slotnum"`
}

func (ship Ship) String() string {
	if len(shipData) > ship.ApiShipId {
		return shipData[ship.ApiShipId].ApiName
	} else {
		return "unknown"
	}
}

var currentDeckId = 1
var decksData [][]Ship

type DeckPort struct {
	ApiShip []int `json:"api_ship"`
}

type ApiPortPort struct {
	ApiBasic    Basic      `json:"api_basic"`
	ApiShip     []Ship     `json:"api_ship"`
	ApiDeckPort []DeckPort `json:"api_deck_port"`
}

func (port ApiPortPort) findShip(id int) (Ship, error) {
	var ship Ship
	for _, ship = range port.ApiShip {
		if ship.ApiId == id {
			return ship, nil
		}
	}
	return ship, errors.New("not found")
}

type KcsapiApiPortPortSimplified struct {
	ApiData []interface{} `json:"api_data"`
	KcsapiBase
}

func handleApiPortPortSimplified(data []byte) error {
	var v KcsapiApiPortPortSimplified
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	fmt.Printf("[%8s] simplified data\n", "Port")
	return err
}

type KcsapiApiPortPort struct {
	ApiData ApiPortPort `json:"api_data"`
	KcsapiBase
}

func handleApiPortPort(data []byte) error {
	var v KcsapiApiPortPort
	err := json.Unmarshal(data, &v)
	if err != nil {
		return handleApiPortPortSimplified(data)
	}

	currentDeckId = 1
	decksData = make([][]Ship, len(v.ApiData.ApiDeckPort))
	lv := (v.ApiData.ApiBasic.ApiLevel / 5)
	if (v.ApiData.ApiBasic.ApiLevel % 5) != 0 {
		lv++
	}
	scouting_lv := float64(lv*5) * 0.6142467
	fmt.Printf("Deck:  1st  2nd  3rd  4th  5th  6th    SKTK (lv=%6.2f)\n", scouting_lv)
	for i, deck := range v.ApiData.ApiDeckPort {
		scouting_ship := 0.0
		scouting_equip := 0.0
		decksData = append(decksData, make([]Ship, 6))
		fmt.Printf("   %d:", i)
		for _, id := range deck.ApiShip {
			ship, err := v.ApiData.findShip(id)
			if err == nil {
				fmt.Printf("%5.0f", ship.ApiCond)
			} else {
				fmt.Printf("  ---")
			}
			decksData[i] = append(decksData[i], ship)

			scouting_ship += math.Sqrt(float64(ship.ApiSakuteki[0])) * 1.6841056
			for _, item := range ship.ApiSlot {
				if item == -1 {
					continue
				}
				switch item {
				case 100:
					fallthrough
				case 200:
					fallthrough
				case 300: /* Bomber */
					scouting_equip += 1 * 1.0376255
				case 500: /* Fighter */
					scouting_equip += 1 * 1.3677954
				case 600: /* Recon Plane */
					scouting_equip += 1 * 1.6592780
				case 700: /* Recon Seaplane */
					scouting_equip += 1 * 2.0000000
				case 750: /* Seaplane Bomber */
					scouting_equip += 1 * 1.7787282
				case 800: /* Small Radar */
					scouting_equip += 1 * 1.0045458
				case 850: /* Large Radar */
					scouting_equip += 1 * 0.9906638
				case 900: /* Searchlight */
					scouting_equip += 1 * 0.9067950
				default:
					/* scouting_equip -= 10000 */
				}
			}
		}
		scouting := scouting_ship + scouting_equip - scouting_lv
		fmt.Printf("%8.2f (ship=%6.2f, equip=%6.2f)\n",
			scouting, scouting_ship, scouting_equip)
	}
	return err
}
