package protocol

import (
	"github.com/tidwall/gjson"
	"fmt"
	"encoding/json"
	"github.com/GreyHood-Studio/server_util/error"
)

type ConnectPacket struct {
	ConnectType	int		`json:ConnectType`
	PlayerName	string	`json:PlayerName`
	Token		string	`json:Token`
}

func PackConnect (connType int, playerName string) []byte{
	fmt.Println(connType, playerName)
	packet := ConnectPacket{ConnectType:connType, PlayerName:playerName, Token:"bitcoin"}
	jsonByte, err := json.Marshal(packet)
	error.CheckError(err, "json pack error")
	return jsonByte
}

func UnpackConnect(data []byte) (int, string, string){
	connectionType := gjson.GetBytes(data, "ConnectType")
	playerName := gjson.GetBytes(data, "PlayerName")
	token := gjson.GetBytes(data, "Token")
	fmt.Println(connectionType, playerName, token)

	return int(connectionType.Int()), playerName.String(), token.String()
}
