package protocol

import (
	"encoding/json"
	"github.com/GreyHood-Studio/play_server/utils"
)

type BulletPacket struct {
	MsgType 	int		`json:msg_type`
	MsgFormat	int		`json:msg_format`
	MsgBody		[]byte	`json:msg_body`
}

func parseBullet(data []byte) BulletPacket {
	var packet BulletPacket
	err := json.Unmarshal(data, &packet)
	utils.CheckError(err, "json unpack error")

	return packet
}