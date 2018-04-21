package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/play_server/utils"
)

type BulletPacket struct {
	PlayerID 	int		`json:"player_id"`
	BulletID	int		`json:"bullet_id"`
	BulletType	int		`json:"bullet_type"`
	BulletPosX	float32	`json:"bulletpos_x"`
	BulletPosY	float32	`json:"bulletpos_y"`
	BulletRotX	float32	`json:"bulletrot_x"`
	BulletRotZ	float32	`json:"bulletrot_z"`
}

func packBullet(packet BulletPacket) []byte {
	jsonByte, err := json.Marshal(packet)
	utils.CheckError(err, "json pack error")
	return jsonByte
}