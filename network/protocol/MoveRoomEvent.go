package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/server_util/checker"
	"github.com/GreyHood-Studio/play_server/model"
)

type MoveRoomEvent struct {
	Address			string				`json:"PlayerNum"`
	Player			model.Player		`json:"MovePlayer"`
}

func PackRoomEvent(packet MoveRoomEvent) []byte{
	jsonByte, err := json.Marshal(packet)
	checker.CheckError(err, "json pack move room event error")

	return jsonByte
}

func UnpackRoomEvent(data []byte) MoveRoomEvent{
	var packet MoveRoomEvent
	err := json.Unmarshal(data, &packet)
	checker.CheckError(err, "unpack player")

	return packet
}
