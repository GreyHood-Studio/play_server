package protocol

type Packet struct {
	Header 	byte		`json:header`
	Ptype	byte		`json:type`
	Pdata	byte		`json:data`
}

// Network에 패킷이 들어왔을 경우, 제일 먼저 처리되는 장소
func RoutePacket() {

}

//
func ParsePacket() {

}