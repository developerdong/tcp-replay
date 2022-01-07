package stream

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"log"
	"net"
	"time"
)

// TcpStreamFactory implements tcpassembly.StreamFactory.
type TcpStreamFactory struct {
	Address  string        // the Address which the copied stream is forwarded to
	Duration time.Duration // how long time waiting for the response from Address
}

func (h *TcpStreamFactory) New(netFlow, transportFlow gopacket.Flow) tcpassembly.Stream {
	conn, err := net.Dial("tcp", h.Address)
	if err != nil {
		log.Fatalln(err)
	}

	tcpConn := conn.(*net.TCPConn)
	if err = tcpConn.CloseRead(); err != nil {
		log.Fatalln(err)
	}

	stream := &TcpStream{
		net:       netFlow,
		transport: transportFlow,
		duration:  h.Duration,
		conn:      tcpConn,
		skip:      false,
	}
	log.Println("Start to copy new stream", netFlow, transportFlow)
	return stream
}
