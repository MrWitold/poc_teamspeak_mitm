package main

import (
	"encoding/binary"
)

/*
SEND:
0:8 - OWN_ID
9:10 - COUNT_BYTES
11:12 - sender_ID
13 - SEPARATOR
14:15 - Message_ID
16 - Codec type
17:n - voice_data
*/

func (user *userClient) preparePacketIn(packet, clientID []byte) []byte{
	tmpCounter := make([]byte, 2)
	packet[0] = user.last13[0]
	packet[1] = user.last13[1]
	packet[2] = user.last13[2]
	packet[3] = user.last13[3]
	packet[4] = user.last13[4]
	packet[5] = user.last13[5]
	packet[6] = user.last13[6]
	packet[7] = user.last13[7]
	//Talk ID
	data := binary.BigEndian.Uint16(user.lastTalkId)
	data++
	binary.BigEndian.PutUint16(tmpCounter, data)
	user.lastTalkId = tmpCounter
	packet[8] = tmpCounter[0]
	packet[9] = tmpCounter[1]

	//uniqe Talk id
	packet[10] = packet[12]
	packet[11] = packet[13]
	packet[12] = packet[14]
	//User talk id
	packet[13] = clientID[1]
	packet[14] = clientID[0]

	return packet
}
/*
RECEIVE:
0:8 - OWN_ID
9:10 - COUNT_BYTES
11 - SEPARATOR
12:13 - Message_ID
14:15 - Sender_ID
16 - Codec type
17:n - voice_data
*/


func (user *userClient) preparePacketOut(packet, clientID []byte) []byte{
	tmpCounter := make([]byte, 2)
	packet[0] = user.last13[0]
	packet[1] = user.last13[1]
	packet[2] = user.last13[2]
	packet[3] = user.last13[3]
	packet[4] = user.last13[4]
	packet[5] = user.last13[5]
	packet[6] = user.last13[6]
	packet[7] = user.last13[7]

	//Talk ID
	data := binary.BigEndian.Uint16(user.lastTalkId)
	data++
	binary.BigEndian.PutUint16(tmpCounter, data)
	user.lastTalkId = tmpCounter
	packet[8] = tmpCounter[0]
	packet[9] = tmpCounter[1]

	//User talk id
	packet[13] = clientID[1]
	packet[14] = clientID[0]

	return packet
}