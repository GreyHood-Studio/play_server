package protocol

import (
	"github.com/tidwall/gjson"
	"fmt"
	"github.com/GreyHood-Studio/play_server/utils"
	"encoding/json"
)

type ConnectPacket struct {
	connectType	int		`json:connectType`
	playerName	[]byte	`json:playerName`
	token		[]byte	`json:token`
}

func PackConnect (connType int, playerName string, token string) []byte{
	event := ConnectPacket{connectType:connType, playerName:[]byte(playerName), token:[]byte(token)}
	jsonByte, err := json.Marshal(event)
	utils.CheckError(err, "json pack error")
	return jsonByte
}

func UnpackConnect(data []byte) (int, string, string){
	connectionType := gjson.GetBytes(data, "connectType")
	playerName := gjson.GetBytes(data, "playerName")
	token := gjson.GetBytes(data, "token")
	fmt.Println(connectionType, playerName, token)

	return int(connectionType.Int()), playerName.String(), token.String()
}
