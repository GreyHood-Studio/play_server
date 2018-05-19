package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/server_util/checker"
)

type CommonEvent struct {
	EventType	int		`json:"EventType"`
	PlayerId	int		`json:"PlayerId"`
	ObjectId	int		`json:"ObjectId,omitempty"`
}

func PackEvent(eventType int, playerId int, objectId int) []byte{
	packet := CommonEvent{EventType:eventType, PlayerId:playerId, ObjectId:objectId}

	jsonByte, err := json.Marshal(packet)
	checker.CheckError(err, "json pack error")
	return jsonByte
}

func UnpackEvent(data []byte) CommonEvent {
	var packet CommonEvent
	err := json.Unmarshal(data, &packet)
	checker.CheckError(err, "unpack common packet")

	return packet
}