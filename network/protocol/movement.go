package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/play_server/utils"
)

type MovementPacket struct {
	PlayerID 	int			`json:"player_id"`
	PlayerPosX	float32		`json:"playerpos_x"`
	PlayerPoxY	float32		`json:"playerpos_y"`
	PlayerRotX	float32		`json:"playerrot_x"`
	PlayerRotY	float32		`json:"playerrot_y"`
}

func parseMovement(data []byte) MovementPacket{
	var packet MovementPacket
	err := json.Unmarshal(data, &packet)
	utils.CheckError(err, "json unpack error")

	return packet
}