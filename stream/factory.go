package stream

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"log"
	"net"
)

// TcpStreamFactory implements tcpassembly.StreamFactory
type TcpStreamFactory struct {
	Address string // the Address which the copied stream is forwarded to
}

func (h *TcpStreamFactory) New(netFlow, transportFlow gopacket.Flow) tcpassembly.Stream {
	conn, err := net.Dial("tcp", h.Address)
	if err != nil {
		log.Fatalln(err)
	}
	stream := &TcpStream{
		net:       netFlow,
		transport: transportFlow,
		r:         tcpreader.NewReaderStream(),
		w:         conn,
	}
	go stream.run() // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &stream.r
}
