package network

import (
	"encoding/json"
	"github.com/GreyHood-Studio/play_server/utils"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/GreyHood-Studio/play_server/network/protocol"
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
		return Packet{-1,-1, []byte{'x'}}
	}
	return packet
}

func packRawPacket(packet Packet) []byte {
	jsonByte, err := json.Marshal(packet)
	utils.CheckError(err, "json pack error")
	return jsonByte
}

func packPacket(mtype int, mfmt int, mbody []byte) []byte {
	packet := Packet {mtype, mfmt, mbody}
	jsonByte, err := json.Marshal(packet)
	utils.CheckError(err, "json pack error")
	fmt.Printf("json marshal %s", jsonByte)
	return jsonByte
}

func routePacket(data []byte) ([]byte, int){
	// unmarshal packet
	var packet Packet
	packet = unpackPacket(data)
	if packet.MsgType == -1 {
		return nil, -1
	}
	// business logic
	result, sendType := processPacket(packet)
	//resultByte := packPacket(result)

	return result, sendType
}

// 처리하는 함수는 TCP Server에서 직접 처리해야하는 것으로 보이며, Protocol은 Packing UnPacking만 담당
// 실제 로직 함수 콜백 함수
func processPacket(packet Packet) ([]byte, int){
	// numbering과 포맷작업 필요
	// gameStart는 response하나, notify하나

	switch packet.MsgFormat {
	case 0: fmt.Println("ping check")
	case 1:
		// redis에서 토큰을 확인하여, 해당 id의 유저가 이 서버의 아이디인지 체크

		return protocol.RequestGameStart(packet.MsgBody), 2
	case 2:
	case 3:
	}

	return []byte{'4'}, 0
}

// 암호화가 있을 경우만 필요
//func decodePacket() {
//
//}