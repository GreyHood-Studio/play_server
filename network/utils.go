package network

import (
	"fmt"
	"encoding/json"
	"github.com/GreyHood-Studio/play_server/network/protocol"
)

func check(err error, message string) {
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", message)
}

func parseData(data []byte) {
	var packet protocol.Packet
	err := json.Unmarshal(data, &packet)
	check(err, "json unpack error")

	//packet.Header
}