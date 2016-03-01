package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

// Build a simple HTTP request parser using tcpassembly.StreamFactory and tcpassembly.Stream interfaces

// HttpStreamFactory implements tcpassembly.StreamFactory
type HttpStreamFactory struct{}

// httpStream will handle the actual decoding of http requests.
type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (h *HttpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hstream.run() // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.r
}

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)
	for {
		switch h.transport.Dst().String() {
		case "80": // request
			r, err := http.ReadRequest(buf)
			if err == io.EOF {
				return
			} else if err != nil {
				log.Println("Error reading stream", h.net, h.transport, ":", err)
			} else {
				r.Body.Close()
				port, _ := strconv.Atoi(h.transport.Src().String())
				reqList.PushBack(Req{port, r.URL.Path})
			}
		default: // response
			r, err := http.ReadResponse(buf, nil)
			if err == io.ErrUnexpectedEOF {
				return
			} else if err != nil {
				log.Println("Error reading stream", h.net, h.transport, ":", err)
			} else {
				var res Res
				res.body, _ = ioutil.ReadAll(r.Body)
				r.Body.Close()
				if r.StatusCode == http.StatusOK {
					res.port, _ = strconv.Atoi(h.transport.Dst().String())
					res.ctype = r.Header["Content-Type"]
					res.cencode = r.Header["Content-Encoding"]
					parserCh <- res
				} else if r.StatusCode >= http.StatusBadRequest {
					log.Println("[Response]:", r.Status)
				}
			}
		}
	}
}
