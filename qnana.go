package main

import (
	"bytes"
	"compress/gzip"
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

var (
	// commandline options
	iface   = flag.String("i", "eth0", "Listen on interface.")
	fname   = flag.String("r", "", "Read packets from file.")
	snaplen = flag.Int("s", 1600, "Snarf bytes of data from each packet.")
	filter  = flag.String("f", "", "Selects which packets will be processed.")

	reqList  = list.New()
	parserCh = make(chan Res)
)

type Req struct {
	port int
	url  string
}

type Res struct {
	port    int
	ctype   []string
	cencode []string
	body    []byte
}

type KcsapiBase struct {
	ApiResult    int    `json:"api_result"`
	ApiResultMsg string `json:"api_result_msg"`
}

type KcsapiGeneral struct {
	ApiData interface{} `json:"api_data"`
	KcsapiBase
}

func handleParseError(url string, data []byte) error {
	var v KcsapiGeneral

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	log.Println("Failed parse JSON:", url)
	str, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("%s\n", str)

	return err
}

func handleGeneral(url string, data []byte) error {
	var v KcsapiGeneral

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	p := "./doc/" + url
	d := path.Dir(p)
	if err := os.MkdirAll(d, 0775); err != nil {
		return err
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	str, _ := json.MarshalIndent(v, "", "  ")
	_, err = f.Write(str)

	return err
}

func parse(wait *sync.WaitGroup) {
	wait.Add(1)

	hasKey := func(a string, list []string) bool {
		for _, b := range list {
			if b == a {
				return true
			}
		}
		return false
	}

	for {
		select {
		case res := <-parserCh:
			if res.port == 0 {
				wait.Done()
				return
			}

			var req Req
			for e := reqList.Front(); e != nil; e = e.Next() {
				if e.Value.(Req).port == res.port {
					req = reqList.Remove(e).(Req)
					break
				}
			}

			if !hasKey("text/plain", res.ctype) {
				continue
			}

			if !strings.HasPrefix(req.url, "/kcsapi/") {
				continue
			}

			var b []byte
			if hasKey("gzip", res.cencode) {
				reader, err := gzip.NewReader(bytes.NewReader(res.body))
				if err != nil {
					log.Println("Failed to parse gzipped data: ", err)
				}
				b, err = ioutil.ReadAll(reader)
				if err != nil {
					log.Println("Failed to read gzipped data: ", err)
				}
			} else {
				b = res.body
			}

			prefix := []byte("svdata=")
			if !bytes.HasPrefix(b, prefix) {
				log.Println("This body is not JSON data")
				continue
			}
			b = bytes.TrimPrefix(b, prefix)

			var err error
			switch req.url {
			case "/kcsapi/api_get_member/base_air_corps":
			case "/kcsapi/api_get_member/basic":
			case "/kcsapi/api_get_member/deck":
			case "/kcsapi/api_get_member/kdock":
			case "/kcsapi/api_get_member/mapcell":
			case "/kcsapi/api_get_member/mapinfo":
			case "/kcsapi/api_get_member/material":
			case "/kcsapi/api_get_member/mission":
			case "/kcsapi/api_get_member/ndock":
			case "/kcsapi/api_get_member/payitem":
			case "/kcsapi/api_get_member/practice":
			case "/kcsapi/api_get_member/preset_deck":
			case "/kcsapi/api_get_member/questlist":
			case "/kcsapi/api_get_member/record":
			case "/kcsapi/api_get_member/require_info":
			case "/kcsapi/api_get_member/ship3":
			case "/kcsapi/api_get_member/ship_deck":
			case "/kcsapi/api_get_member/slot_item":
			case "/kcsapi/api_get_member/sortie_conditions":
			case "/kcsapi/api_get_member/unsetslot":
			case "/kcsapi/api_get_member/useitem":
				/* do nothing */
			case "/kcsapi/api_port/port":
				err = handleApiPortPort(b)
			case "/kcsapi/api_req_member/get_practice_enemyinfo":
				/* do nothing */
			case "/kcsapi/api_req_map/start":
				err = handleApiReqMapStart(b)
			case "/kcsapi/api_req_map/next":
				err = handleApiReqMapNext(b)
			case "/kcsapi/api_req_combined_battle/battleresult":
				err = handleApiReqCombinedBattleBattleresult(b)
			case "/kcsapi/api_req_combined_battle/battle_water":
				err = handleApiReqCombinedBattleBattleWater(b)
			case "/kcsapi/api_req_combined_battle/goback_port":
				/* do nothing */
			case "/kcsapi/api_req_combined_battle/battle":
				err = handleApiReqCombinedBattleBattle(b)
			case "/kcsapi/api_req_combined_battle/ld_airbattle":
				err = handleApiReqCombinedBattleLdAirbattle(b)
			case "/kcsapi/api_req_combined_battle/midnight_battle":
				err = handleApiReqCombinedBattleMidnightBattle(b)
			case "/kcsapi/api_req_practice/battle":
				err = handleApiReqPracticeBattle(b)
			case "/kcsapi/api_req_practice/midnight_battle":
				err = handleApiReqPracticeMidnightBattle(b)
			case "/kcsapi/api_req_practice/battle_result":
				err = handleApiReqPracticeBattleResult(b)
			case "/kcsapi/api_req_sortie/battle":
				err = handleApiReqSortieBattle(b)
			case "/kcsapi/api_req_sortie/battleresult":
				err = handleApiReqSortieBattleresult(b)
			case "/kcsapi/api_req_sortie/airbattle":
				err = handleApiReqSortieAirbattle(b)
			case "/kcsapi/api_req_sortie/ld_airbattle":
				err = handleApiReqSortieLdAirbattle(b)
			case "/kcsapi/api_req_battle_midnight/battle":
				err = handleApiReqBattleMidnightBattle(b)
			case "/kcsapi/api_req_battle_midnight/sp_midnight":
				err = handleApiReqBattleMidnightBattle(b)
			case "/kcsapi/api_req_air_corps/set_action":
			case "/kcsapi/api_req_air_corps/set_plane":
			case "/kcsapi/api_req_air_corps/supply":
			case "/kcsapi/api_req_furniture/buy":
			case "/kcsapi/api_req_furniture/change":
			case "/kcsapi/api_req_hensei/change":
			case "/kcsapi/api_req_hensei/combined":
			case "/kcsapi/api_req_hensei/lock":
			case "/kcsapi/api_req_hensei/preset_register":
			case "/kcsapi/api_req_hensei/preset_select":
			case "/kcsapi/api_req_hokyu/charge":
			case "/kcsapi/api_req_kaisou/lock":
			case "/kcsapi/api_req_kaisou/powerup":
			case "/kcsapi/api_req_kaisou/remodeling":
			case "/kcsapi/api_req_kaisou/slot_exchange_index":
			case "/kcsapi/api_req_kaisou/slotset":
			case "/kcsapi/api_req_kaisou/unsetslot_all":
			case "/kcsapi/api_req_kousyou/createitem":
			case "/kcsapi/api_req_kousyou/createship":
			case "/kcsapi/api_req_kousyou/destroyitem2":
			case "/kcsapi/api_req_kousyou/destroyship":
			case "/kcsapi/api_req_kousyou/getship":
			case "/kcsapi/api_req_kousyou/remodel_slot":
			case "/kcsapi/api_req_kousyou/remodel_slotlist":
			case "/kcsapi/api_req_kousyou/remodel_slotlist_detail":
			case "/kcsapi/api_req_map/select_eventmap_rank":
			case "/kcsapi/api_req_map/start_air_base":
			case "/kcsapi/api_req_member/get_incentive":
			case "/kcsapi/api_req_member/itemuse":
			case "/kcsapi/api_req_mission/result":
			case "/kcsapi/api_req_mission/start":
			case "/kcsapi/api_req_nyukyo/start":
			case "/kcsapi/api_req_nyukyo/speedchange":
			case "/kcsapi/api_req_quest/clearitemget":
			case "/kcsapi/api_req_quest/start":
			case "/kcsapi/api_req_quest/stop":
				/* do nothing */
			case "/kcsapi/api_start2":
				/* allways update data */
				err = handleGeneral(req.url, b)
				if err != nil {
					err = handleApiStart2(b)
				}
			default:
				log.Println("Unknown API:", req.url)
				err = handleGeneral(req.url, b)
			}
			if err != nil {
				handleParseError(req.url, b)
			}
		}
	}
}

func main() {
	var handle *pcap.Handle
	var err error

	flag.Parse()

	if *fname != "" {
		log.Printf("Reading from pcap file %q", *fname)
		handle, err = pcap.OpenOffline(*fname)
	} else {
		log.Printf("Starting capture on interface %q", *iface)
		handle, err = pcap.OpenLive(*iface, int32(*snaplen), false, pcap.BlockForever)
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := dataSet.loadLocal("./doc/kcsapi/api_start2"); err != nil {
		log.Println("Failed to load local data set: ", err)
	}

	if *filter == "" {
		// from kcs_const.js
		*filter = "tcp port 80 and host"
		*filter += " (125.6.184.15 or 125.6.184.16 or 125.6.187.205 or"
		*filter += " 125.6.187.229 or 125.6.187.253 or 125.6.188.25 or"
		*filter += " 125.6.189.7 or 125.6.189.39 or 125.6.189.71 or"
		*filter += " 125.6.189.103 or 125.6.189.135 or 125.6.189.167 or"
		*filter += " 125.6.189.215 or 125.6.189.247 or 203.104.209.7 or"
		*filter += " 203.104.209.23 or 203.104.209.39 or 203.104.209.55 or"
		*filter += " 203.104.209.71 or 203.104.209.102 or 203.104.248.135)"
	}
	if err := handle.SetBPFFilter(*filter); err != nil {
		log.Fatal(err)
	}

	// Set up assembly
	streamFactory := &HttpStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	// Set up paser goroutine
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	wait := new(sync.WaitGroup)
	go parse(wait)

	// Read in packets, pass to assembler.
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates the end of a pcap file.
			if packet == nil {
				return
			}

			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil {
				continue
			}
			if packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				continue
			}

			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

		case <-ticker:
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}

	parserCh <- Res{port: 0}
	wait.Wait()
}
