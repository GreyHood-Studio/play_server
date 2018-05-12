package protocol

import (
	"github.com/tidwall/sjson"
	"github.com/GreyHood-Studio/server_util/error"
)

type BulletPacket struct {
	PlayerID	int		`json:PlayerID`
	BulletID	int		`json:BulletID`
	FirePosX	float32	`json:FirePosX`
	FirePosY	float32	`json:FirePosY`
	MousePosX	float32	`json:MousePosX`
	MousePosY	float32	`json:MousePosY`
}

func AssignBulletID(bulletID int, data []byte) []byte{
	reset, err := sjson.SetBytes(data, "BulletID", bulletID)
	error.NoDeadError(err, "assign bullet id error")
	return reset
}
