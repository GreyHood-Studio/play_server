package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/play_server/utils"
)

type PlayerInfoPacket struct {
	PlayerId		int		`json:"player_id"`
	PlayerHealth	int		`json:"player_health"`
}

func parsePlayerInfo(data []byte) PlayerInfoPacket{
	var packet PlayerInfoPacket
	err := json.Unmarshal(data, &packet)
	utils.CheckError(err, "json unpack error")

	return packet
}