package main

import (
	"encoding/json"
	"fmt"
)

type Ship struct {
	ApiId int			`json:"api_id"`
	ApiCond float64		`json:"api_cond"`
}

type DeckPort struct {
	ApiShip []int		`json:"api_ship"`
}

type ApiPortPort struct {
	ApiShip []Ship			`json:"api_ship"`
	ApiDeckPort []DeckPort	`json:"api_deck_port"`
}

type KcsapiApiPortPort struct {
	ApiData ApiPortPort		`json:"api_data"`
	KcsapiBase
}

func handleApiPortPort(data []byte) error {
	var v KcsapiApiPortPort
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}


	fmt.Printf("Deck:  1st  2nd  3rd  4th  5th  6th\n")
	for i, deck := range v.ApiData.ApiDeckPort {
		fmt.Printf("   %d:", i)
		for _, ship := range deck.ApiShip {
			for _, s := range v.ApiData.ApiShip {
				if ship == s.ApiId {
					fmt.Printf("%5.0f", s.ApiCond)
					break
				}
			}
		}
		fmt.Printf("\n")
	}
	return err
}

