/*
Package stream builds a simple TCP parser using tcpassembly.StreamFactory and tcpassembly.Stream interfaces
*/
package stream

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"io"
	"log"
	"net"
	"sync"
)

// TcpStream will handle the actual forwarding of copied tcp streams.
type TcpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
	c              *net.TCPConn
}

func (t *TcpStream) run() {
	log.Println("Start to copy new stream", t.net, t.transport)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		// discard the response
		_, _ = io.Copy(io.Discard, t.c)
		_ = t.c.CloseRead()
		wg.Done()
	}()
	go func() {
		// forward copied data to the remote address
		if _, err := io.Copy(t.c, &t.r); err != nil {
			log.Println("Error reading stream", t.net, t.transport, ":", err)
		}
		_ = t.c.CloseWrite()
		wg.Done()
	}()
	wg.Wait()
	log.Println("Finish copying new stream", t.net, t.transport)
}
