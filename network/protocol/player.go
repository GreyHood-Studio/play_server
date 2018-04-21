package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/play_server/utils"
)

type PlayerInfoPacket struct {
	PlayerId		int		`json:"player_id"`
	PlayerHealth	int		`json:"player_health"`
}

func packPlayerInfo(packet PlayerInfoPacket) []byte {
	jsonByte, err := json.Marshal(packet)
	utils.CheckError(err, "json pack error")
	return jsonByte
}