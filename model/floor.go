package model

// 한 층의 데이터 일반적인 게임에서의 room 개념과 유사
type Floor struct {
	id			uint32
	status		FloorStatus	// 맵에 대한 현재 상태 구조체
	characters	[]Player	// 맵에 존재하는 플레이어들의 리스트 ( user, enemy, character )
	record		Record		// 이동, 피격 판정에 대한 데이터 큐
}

// 한 층의 현재 상태 ( 동접, 방 관리를 위한 내용 )
// 플레이어 수, ai수, etc...
type FloorStatus struct {
	port			int16
	playerCount		int
}

// 피격 판정과 이동 판정을 위한 기록 데이터
type Record struct {

}