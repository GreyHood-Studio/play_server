package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/server_util/checker"
	"github.com/GreyHood-Studio/play_server/model"
	"github.com/tidwall/gjson"
)

type Item struct {
	ItemId 				int 				`json:"ItemId"`
}

func ExtractPlayerInStartPacket(data []byte) (*model.Player) {
	//mapId := gjson.GetBytes(data, "MapId")
	//spawnId := gjson.GetBytes(data, "SpawnId")
	playerName := gjson.GetBytes(data, "ConnectInfos.0.PlayerName")
	//playerId := gjson.GetBytes(data, "ConnectInfos.0.PlayerId")
	currentHealth := gjson.GetBytes(data, "ConnectInfos.0.Health")
	currentWeaponId := gjson.GetBytes(data, "ConnectInfos.0.CurrWeaponId")

	newPlayer := &model.Player{
		PlayerName: playerName.String(), CurrentHealth: int(currentHealth.Int()), CurrentWeaponId: int(currentWeaponId.Int()) }

	//packet := StartEvent{ MapId:int(mapId.Int()), SpawnId: spawnId}
	//packet.PlayerList = append(packet.PlayerList, playerEvent)
	return newPlayer
}

func PackPlayer (model model.Player) []byte{
	jsonByte, err := json.Marshal(model)
	checker.CheckError(err, "json pack checker")
	return jsonByte
}

func UnpackPlayer(data []byte) model.Player {
	var packet model.Player
	err := json.Unmarshal(data, &packet)
	checker.CheckError(err, "unpack player")
	return packet
}