package protocol

import (
	"github.com/tidwall/gjson"
	"encoding/json"
	"github.com/GreyHood-Studio/server_util/error"
	"fmt"
)

type PlayerEvent struct {
	PlayerName	string	`json:PlayerName`
	PlayerId	int		`json:PlayerId`
	SpawnId		int		`json:SpawnId`
	Items		int		`json:Items`
}

func PackPlayer (event PlayerEvent) []byte{
	jsonByte, err := json.Marshal(event)
	error.CheckError(err, "json pack error")
	return jsonByte
}

func UnpackPlayer(data []byte) PlayerEvent {
	playerName := gjson.GetBytes(data, "PlayerName")
	playerId := gjson.GetBytes(data, "PlayerId")
	spawnId := gjson.GetBytes(data, "SpawnId")
	fmt.Println(playerName, playerId, spawnId)

	packet := PlayerEvent{ PlayerName:playerName.String(), PlayerId:int(playerId.Int()), SpawnId:int(spawnId.Int()) }
	return packet
}
