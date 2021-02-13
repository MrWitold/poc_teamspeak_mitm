package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"log"
	"os"
	"time"
)

func (proxy *proxyInterface) readFromFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r, err := pcapgo.NewNgReader(f, pcapgo.DefaultNgReaderOptions)
	if err != nil {
		log.Fatal(err)
	}

	for {
		data, _, err := r.ReadPacketData()
		if err != nil {
			if err.Error() == "EOF"{
				fmt.Println("END OF FILE")
				return
			}
			log.Fatal(err)
		}
		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
		if packet == nil {
			fmt.Println("err 1")
			return
		}

		/*
			if tcpLayer := packet.Layer(layers.LayerTypeUDP); tcpLayer != nil {
				// Get actual TCP data from this layer
				udp, _ := tcpLayer.(*layers.UDP)

				if udp.SrcPort != 9987 {
					continue
				}
				if udp.DstPort == 59314 {
					continue
				}
			}
		*/

		netLayer := packet.NetworkLayer()
		if nil == netLayer {
			fmt.Println("err 2")
			continue
		}
		newFlow := netLayer.NetworkFlow()
		_, dst := newFlow.Endpoints()

		if app := packet.ApplicationLayer(); app != nil {
			n := len(app.Payload())
			if n < 15 {
				continue
			}

			if dst.String() == "91.299.227.192"{
				proxy.usersMap[proxy.attacker].injectPacket <- capturePacket{data: app.Payload()[0:n],typeOfPacket: "out"}
				time.Sleep(11*time.Millisecond)
			}
			/*
				if src.String() == "91.299.227.192" {
					proxy.usersMap[proxy.attacker].injectPacket <- capturePacket{data:app.Payload()[0:n],typeOfPacket: "in"}
					time.Sleep(11*time.Millisecond)
				}
			*/
			//fmt.Println(app.Payload())
		}
	}
}