package main

import (
	"fmt"
)

type Hps struct {
	now   []int
	max   []int
	dmg   []int
	label string
}

func (h Hps) dumpHps() {
	fmt.Printf("[%8s]:", h.label)
	for i, hp := range h.now {
		var flag string
		hp -= h.dmg[i]
		if h.max[i] < 0 {
			if i != 0 {
				fmt.Printf(" ---/--- ")
			}
			continue
		} else if hp > (h.max[i] * 3 / 4) {
			flag = " "
		} else if hp > (h.max[i] / 2) {
			flag = "-"
		} else if hp > (h.max[i] / 4) {
			flag = "="
		} else {
			flag = "!"
		}
		fmt.Printf(" %3d/%3d%s", hp, h.max[i], flag)
	}
	fmt.Printf("\n")
}

type Damage struct {
	deck  []Hps
	enemy Hps
}

func (d *Damage) init(now []int, max []int, enemy []int) {
	d.deck = make([]Hps, 2)
	enemy_size := len(enemy) - 1
	d.deck[0].now = now[:len(now)-enemy_size]
	d.deck[0].max = max[:len(max)-enemy_size]
	d.deck[0].dmg = make([]int, len(now)-enemy_size)
	d.deck[0].label = "Deck"
	d.enemy.now = enemy[len(now)-enemy_size:]
	d.enemy.max = enemy[len(max)-enemy_size:]
	d.enemy.dmg = make([]int, enemy_size)
	d.enemy.label = "Enemy"
	d.deck[1].now = nil
}

func (d *Damage) initCombined(now []int, max []int) {
	d.deck[1].now = now
	d.deck[1].max = max
	d.deck[1].dmg = make([]int, len(now))
	d.deck[0].label = "Deck1"
	d.deck[1].label = "Deck2"
}

func (d Damage) dumpHps() {
	d.deck[0].dumpHps()
	if d.deck[1].now != nil {
		d.deck[1].dumpHps()
	}
	d.enemy.dumpHps()
}

func dumpHps(label string, hps []int, maxhps []int) {
	fmt.Printf("[%8s]:", label)
	for i, hp := range hps {
		var flag string
		if maxhps[i+1] < 0 {
			fmt.Printf(" ---/--- ")
			continue
		} else if hp > (maxhps[i+1] * 3 / 4) {
			flag = " "
		} else if hp > (maxhps[i+1] / 2) {
			flag = "-"
		} else if hp > (maxhps[i+1] / 4) {
			flag = "="
		} else {
			flag = "!"
		}
		fmt.Printf(" %3d/%3d%s", hp, maxhps[i+1], flag)
	}
	fmt.Printf("\n")
}
