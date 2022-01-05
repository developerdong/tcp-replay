/*
Package stream builds a simple TCP parser using tcpassembly.StreamFactory and tcpassembly.Stream interfaces
*/
package stream

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net"
	"time"
)

// TcpStream will handle the actual forwarding of copied tcp streams.
type TcpStream struct {
	net, transport gopacket.Flow
	duration       time.Duration
	r              tcpreader.ReaderStream
	c              *net.TCPConn
}

func (t *TcpStream) run() {
	log.Println("Start to copy new stream", t.net, t.transport)
	g := new(errgroup.Group)
	g.Go(func() error {
		// discard the response
		_, err := io.Copy(io.Discard, t.c)
		_ = t.c.CloseRead()
		return err
	})
	g.Go(func() error {
		// forward copied data to the remote address
		_, err := io.Copy(t.c, &t.r)
		_ = t.r.Close()
		_ = t.c.CloseWrite()
		return err
	})
	if err := g.Wait(); err != nil {
		log.Println("Error reading stream", t.net, t.transport, ":", err)
	} else {
		// If the original connection terminates correctly, wait for the response from
		// new connection. Otherwise, the request in the new connection may be cancelled.
		time.Sleep(t.duration)
	}
	log.Println("Finish copying new stream", t.net, t.transport)
}
