package protocol

import (
	"github.com/tidwall/sjson"
	"github.com/GreyHood-Studio/server_util/checker"
	"github.com/tidwall/gjson"
)

type BulletPacket struct {
	PlayerNum	int			`json:"PlayerNum"`
	PositionX	float32		`json:"PositionX"`
	PositionY	float32		`json:"PositionY"`
	MouseX		float32		`json:"MouseX"`
	MouseY		float32		`json:"MouseY"`

	WeaponId	int			`json:"WeaponId"`
	BulletID	int			`json:"BulletId"`
	BulletNum	int			`json:"BulletNum"`
}

func CalculateDamage(weaponId int, bulletId int) int {
	return 1
}

func AssignBulletNum(bulletNum int, data []byte) []byte{
	reset, err := sjson.SetBytes(data, "BulletNum", bulletNum)
	checker.NoDeadError(err, "assign bullet id checker")
	return reset
}

func ExtractBulletType(data []byte) (int){
	weaponId := gjson.GetBytes(data, "WeaponId")
	bulletId := gjson.GetBytes(data, "BulletId")

	bulletDamage := CalculateDamage(int(weaponId.Int()), int(bulletId.Int()))
	return bulletDamage
}
