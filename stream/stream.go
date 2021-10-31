/*
Package stream builds a simple TCP parser using tcpassembly.StreamFactory and tcpassembly.Stream interfaces
*/
package stream

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"io"
	"log"
)

// TcpStream will handle the actual forwarding of copied tcp streams.
type TcpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
	w              io.WriteCloser
}

func (t *TcpStream) run() {
	if _, err := io.Copy(t.w, &t.r); err != nil {
		log.Println("Error reading stream", t.net, t.transport, ":", err)
	}
	_ = t.w.Close()
}
