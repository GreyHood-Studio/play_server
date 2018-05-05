package protocol

import (
	"github.com/tidwall/gjson"
	"fmt"
)

type PlayerEvent struct {
	PlayerName	string	`json:PlayerName`
	PlayerId	int		`json:PlayerId`
	SpawnId		int		`json:SpawnId`
	Items		int		`json:Items`
}

func UnpackPlayer(data []byte) PlayerEvent {
	playerName := gjson.GetBytes(data, "PlayerName")
	playerId := gjson.GetBytes(data, "PlayerId")
	spawnId := gjson.GetBytes(data, "SpawnId")
	fmt.Println(playerName, playerId, spawnId)

	packet := PlayerEvent{ PlayerName:playerName.String(), PlayerId:int(playerId.Int()), SpawnId:int(spawnId.Int()) }
	return packet
}
