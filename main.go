/*
	This binary reads packets off the wire and reconstructs TCP content it sees, forwarding them.
*/
package main

import (
	"flag"
	"github.com/developerdong/tcp-replay/stream"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"log"
	"os"
	"time"
)

func main() {
	var inputFile = flag.String("i", "", "Filename to read from, stdin is the default")
	var filter = flag.String("f", "", "BPF filter for pcap")
	var logAllPackets = flag.Bool("v", false, "Logs every packet in great detail")
	var targetAddress = flag.String("t", "localhost:8080", "Target address a copied stream is forwarded to")
	var duration = flag.Duration("d", 0, "how long time waiting for the response from target address after the original connection terminating normally")
	flag.Parse()

	var handle *pcap.Handle
	var err error

	// Set up pcap packet capture
	if *inputFile != "" {
		handle, err = pcap.OpenOffline(*inputFile)
	} else {
		handle, err = pcap.OpenOfflineFile(os.Stdin)
	}
	if err != nil {
		log.Fatal(err)
	}

	if err = handle.SetBPFFilter(*filter); err != nil {
		log.Fatal(err)
	}

	// Set up assembly
	streamFactory := &stream.TcpStreamFactory{Address: *targetAddress, Duration: *duration}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	log.Println("reading in packets")
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
			if *logAllPackets {
				log.Println(packet)
			}
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Println("Unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

		case <-ticker:
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}
