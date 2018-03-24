package model

// 플레이어 객체 정보
type Player struct {
	currentPos		Pos
	currentFloor 	int
	currentRoom		int
}

type Enemy struct {
	enemyType		int
	enemySize		int
}

type User struct {
	sessionid		string
	userName		string
	playinfo		Player
}