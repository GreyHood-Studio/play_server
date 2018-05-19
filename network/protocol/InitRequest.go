package protocol

import (
	"fmt"
	"encoding/json"
	"github.com/GreyHood-Studio/server_util/checker"
	"bytes"
)

type InitRequest struct {
	ConnectType	int		`json:"ConnectType"`
	Token		string	`json:"Token,omitempty"`
	PlayerName	string	`json:"PlayerName"`
}

func PackInitPacket (connType int, playerName string) []byte{
	packet := InitRequest{ConnectType:connType, PlayerName:playerName, Token:"default"}
	jsonByte, err := json.Marshal(packet)
	checker.CheckError(err, "json pack checker")
	return jsonByte
}

func UnpackInitPacket(data []byte) InitRequest{
	var packet InitRequest

	trimBytes := bytes.TrimRight(data, "\x00")

	fmt.Println("unpack init", string(trimBytes))
	err := json.Unmarshal(trimBytes, &packet)

	checker.CheckError(err, "unpack init packet")

	return packet
}
