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
	"runtime"
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

func handleGeneral(data []byte) error {
	var v KcsapiGeneral

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

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
			case "/kcsapi/api_port/port":
				err = handleApiPortPort(b)
			case "/kcsapi/api_req_map/next":
				/* do nothing */
			case "/kcsapi/api_get_member/ship_deck":
				/* do nothing */
			case "/kcsapi/api_req_combined_battle/battle_water":
				err = handleApiReqCombinedBattleBattleWater(b)
			case "/kcsapi/api_req_combined_battle/ld_airbattle":
				err = handleApiReqCombinedBattleLdAirbattle(b)
			case "/kcsapi/api_req_combined_battle/battleresult":
				err = handleApiReqCombinedBattleBattleresult(b)
			case "/kcsapi/api_req_practice/battle":
				err = handleApiReqPracticeBattle(b)
			case "/kcsapi/api_req_practice/midnight_battle":
				err = handleApiReqPracticeMidnightBattle(b)
			case "/kcsapi/api_req_practice/battle_result":
				err = handleApiReqPracticeBattleResult(b)
			case "/kcsapi/api_req_sortie/battle":
				err = handleApiReqSortieBattle(b)
			case "/kcsapi/api_req_battle_midnight/battle":
				handleParseError(req.url, b)
				err = handleApiReqBattleMidnightBattle(b)
			case "/kcsapi/api_req_sortie/battleresult":
				err = handleApiReqSortieBattleresult(b)
			default:
				log.Println("Unknown API:", req.url)
				err = handleGeneral(b)
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
