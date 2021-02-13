package main

import (
	"fmt"
	"net"
	_ "net/http/pprof"
	"time"
)

type proxyInterface struct {
	clientToProxy    chan *packet
	startCapture     bool
	teamspeakAddress string
	usersMap         map[string]*userClient
	attacker         string
	victim           string
	outputUser       uint16
	endSignal        chan int
}

type userClient struct {
	proxy                  *proxyInterface
	injectPacket           chan capturePacket
	incomingPackets        chan []byte
	last13                 []byte
	lastTalkId             []byte
	userAddressIP          string
	userTeamspeakAddressIP string
	end                    chan int
}

type capturePacket struct {
	data         []byte
	typeOfPacket string
}

type packet struct {
	data    []byte
	address string
}

func (proxy *proxyInterface) setupProxy(address string) {
	s, _ := net.ResolveUDPAddr("udp4", address)
	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The UDP server is %v\n", s)
	defer connection.Close()

	go proxy.readInputPacket(connection)
	go proxy.sendOutputPacket(connection)

	<-proxy.endSignal
}

func (user *userClient) UserTeamspeakConnection(address string) {
	s, _ := net.ResolveUDPAddr("udp4", address)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	user.userTeamspeakAddressIP = c.LocalAddr().String()

	fmt.Printf("%v :The UDP client is %v -> %s \n", time.Now().Format("2006-01-02 15:04:05"), user.userTeamspeakAddressIP, user.userAddressIP)
	defer c.Close()

	go user.sendPacketToTeamspeak(c)
	go user.receivePacketFromTeamspeak(c)

	<-user.end
}

func main() {
	proxyClient := &proxyInterface{clientToProxy: make(chan *packet, 4096),
		startCapture:     false,
		teamspeakAddress: "",
		attacker:         "",
		victim:           "",
		outputUser:       uint16(0),
		usersMap:         make(map[string]*userClient),
	}
	go proxyClient.setupProxy(":324")

	for {
		proxyClient.setupMenu()
	}
}
