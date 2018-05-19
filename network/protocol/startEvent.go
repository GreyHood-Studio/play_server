package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/server_util/checker"
	"github.com/GreyHood-Studio/play_server/model"
)

type StartEvent struct {
	MapId			int				`json:"MapId"`
	SpawnId			int				`json:"SpawnId"`
	ConnectionInfo	[]*model.Player	`json:"ConnectInfos"`
}

func PackStart(event StartEvent) []byte{
	jsonByte, err := json.Marshal(event)
	checker.CheckError(err, "json pack error")
	return jsonByte
}

func NewStartPacket(mapId int, spawnId int, playerCount int) *StartEvent{
	startEvent := &StartEvent{
		MapId: mapId,
		SpawnId: spawnId,
		ConnectionInfo: make([]*model.Player, playerCount),
	}

	return startEvent
}