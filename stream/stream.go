package stream

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"log"
	"net"
	"time"
)

// TcpStream will handle the actual forwarding of copied tcp streams.
type TcpStream struct {
	net, transport gopacket.Flow
	duration       time.Duration
	conn           *net.TCPConn
	skip           bool
}

// Reassembled implements tcpassembly.Stream's Reassembled function.
func (t *TcpStream) Reassembled(reassembly []tcpassembly.Reassembly) {
	for _, r := range reassembly {
		if r.Skip != 0 {
			// We didn't capture the whole data.
			log.Println("Data is skipped", t.net, t.transport, ":", r.Skip)
			t.skip = true
		}
		if t.skip {
			// We can't forward data anymore if some bytes were skipped.
			return
		}
		if _, err := t.conn.Write(r.Bytes); err != nil {
			log.Println("Error processing stream", t.net, t.transport, ":", err)
			t.skip = true
		}
	}
}

// ReassemblyComplete implements tcpassembly.Stream's ReassemblyComplete
// function.
func (t *TcpStream) ReassemblyComplete() {
	go func() {
		if !t.skip {
			// If we have captured the whole data from the original connection, wait for the
			// response from new connection. Otherwise, the request in the new connection may
			// be cancelled.
			time.Sleep(t.duration)
		}
		_ = t.conn.Close()
		log.Println("Finish copying new stream", t.net, t.transport)
	}()
}
