package protocol

// json으로 처리해야함
type Protocol struct{
	msg_type	string
	msg_data	string
}

// Network에 패킷이 들어왔을 경우, 제일 먼저 처리되는 장소
func RoutePacket() {

}

//
func ParsePacket() {

}