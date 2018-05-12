package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/server_util/error"
	"github.com/tidwall/gjson"
)

type StartEvent struct {
	MapId		int				`json:MapId`
	PlayerList	[]PlayerEvent	`json:PlayerList`
	HistoryList []CommonEvent	`json:HistoryList`
}

func PackStart(event StartEvent) []byte{
	dumpEvent := CommonEvent{EventType:0, PlayerId:0, ObjectId:0}
	event.HistoryList = append(event.HistoryList, dumpEvent)
	jsonByte, err := json.Marshal(event)
	error.CheckError(err, "json pack error")
	return jsonByte
}

func UnpackStart(data []byte) StartEvent {
	mapId := gjson.GetBytes(data, "MapId")
	playerName := gjson.GetBytes(data, "PlayerList.0.player_name")
	playerId := gjson.GetBytes(data, "PlayerList.0.player_id")
	spawnId := gjson.GetBytes(data, "PlayerList.0.spawn_position")

	playerEvent := PlayerEvent{
		PlayerName: playerName.String(), PlayerId:int(playerId.Int()), SpawnId:int(spawnId.Int()) }

	packet := StartEvent{ MapId:int(mapId.Int())}
	packet.PlayerList = append(packet.PlayerList, playerEvent)
	return packet
}