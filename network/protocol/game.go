package protocol

import "github.com/tidwall/gjson"

func RequestGameStart(data []byte) []byte {
	playerId := gjson.GetBytes(data, "player_id")

	// 레디스로 해당 유저가 해당 서버에 오는게 맞는지 체크
	if validationUser(playerId.String()) {

	}

	return []byte{'d'}
}

func validationUser(playerId string) bool {

	return true
}