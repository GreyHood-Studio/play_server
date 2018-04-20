package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/play_server/utils"
)

type Packet struct {
	MsgType 	int		`json:"msg_type"`
	MsgFormat	int		`json:"msg_format"`
	MsgBody		[]byte	`json:"msg_body"`
}

func ParsePacket(data []byte) Packet {
	var packet Packet
	err := json.Unmarshal(data, &packet)
	utils.CheckError(err, "json unpack error")

	return packet
}

func RoutePacket(data []byte) ([]byte, int){
	packet := ParsePacket(data)

	// 이제 비지니스 로직에 따라 처리해서 데이터를 가져오는 작업이 필요
	result, err := json.Marshal(packet)
	utils.CheckError(err, "json encoding error")
	return result, 0
}

// 실제 로직 함수 콜백 함수
func processPacket(packet Packet) {
	switch packet.MsgFormat {
	case 0:
	case 1:
	case 2:
	case 3:
	}
}

// 암호화가 있을 경우만 필요
func decodePacket() {

}