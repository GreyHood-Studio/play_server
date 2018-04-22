package network

import (
	"encoding/json"
	"github.com/GreyHood-Studio/play_server/utils"
	"fmt"
	"github.com/tidwall/gjson"
)

type Packet struct {
	MsgType 	int		`json:"msg_type"`
	MsgFormat	int		`json:"msg_format"`
	MsgBody		[]byte	`json:"msg_body"`
}

func unpackPacket(data []byte) Packet {
	var packet Packet
	test := gjson.GetBytes(data, "msg_type")
	fmt.Println(test)
	packet, ok := gjson.ParseBytes(data).Value().(Packet)
	if !ok {
		fmt.Println("json parse error")
		return Packet{-2,-1, []byte{'x'}}
	}
	return packet
}

func packRawPacket(packet Packet) ([]byte, int) {
	jsonByte, err := json.Marshal(packet)
	utils.CheckError(err, "json pack error")
	return jsonByte, packet.MsgType
}

func packValuePacket(mtype int, mfmt int, mbody []byte) []byte {
	packet := Packet {mtype, mfmt, mbody}
	jsonByte, err := json.Marshal(packet)
	utils.CheckError(err, "json pack error")
	fmt.Printf("json marshal %s", jsonByte)
	return jsonByte
}

// 암호화가 있을 경우만 필요
//func decodePacket() {
//
//}