package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/play_server/utils"
)

type GameStartRequestPacket struct {
	PlayerID 	int			`json:"player_id"`
	StartPosX	float32		`json:"startpos_x"`
	StartPoxY	float32		`json:"startpos_y"`
}

type GameStartNotifyPacket struct {
	PlayerID 	int			`json:"player_id"`
	StartPosX	float32		`json:"startpos_x"`
	StartPoxY	float32		`json:"startpos_y"`
}

func parseGameStart(data []byte) GameStartRequestPacket {
	var packet GameStartRequestPacket
	err := json.Unmarshal(data, &packet)
	utils.CheckError(err, "json unpack error")

	return packet
}

func requestGameStart() {

}