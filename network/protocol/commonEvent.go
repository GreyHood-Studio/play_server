package protocol

import (
	"github.com/tidwall/gjson"
	"fmt"
	"encoding/json"
	"github.com/GreyHood-Studio/server_util/error"
)

type CommonEvent struct {
	EventType	int		`json:EventType`
	PlayerId	int		`json:PlayerId`
	ObjectId	int		`json:ObjectId`
}

func PackEvent(eventType int, playerId int, objectId int) []byte{
	packet := CommonEvent{EventType:eventType, PlayerId:playerId, ObjectId:objectId}

	jsonByte, err := json.Marshal(packet)
	error.CheckError(err, "json pack error")
	return jsonByte
}

func UnpackEvent(data []byte) CommonEvent {
	eventType := gjson.GetBytes(data, "EventType")
	playerId := gjson.GetBytes(data, "PlayerId")
	objectId := gjson.GetBytes(data, "ObjectId")
	fmt.Println(eventType, playerId, objectId)

	connectEvent := CommonEvent{ EventType:int(eventType.Int()),
		PlayerId:int(playerId.Int()), ObjectId:int(objectId.Int()) }
	return connectEvent
}