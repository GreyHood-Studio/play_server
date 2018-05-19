package protocol

import (
	"github.com/GreyHood-Studio/server_util/checker"
	"encoding/json"
)

type HitPacket struct {
	// event type이 1일때는 단순 피격, 2일때는 사망
	// 데미지 계산은 서버가 알아서
	EventType		int			`json:"EventType"`
	PlayerNum		int			`json:"PlayerNum"`
	BulletNum		int			`json:"BulletNum"`
}

func PackHit (event HitPacket) []byte{
	jsonByte, err := json.Marshal(event)
	checker.CheckError(err, "json pack checker")
	return jsonByte
}

func UnpackHit (data []byte) HitPacket{
	var packet HitPacket
	err := json.Unmarshal(data, &packet)
	checker.CheckError(err, "unpack hit")
	return packet
}