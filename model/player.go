package model

// 플레이어 객체 정보
type Player struct {
	PlayerNum			int					`json:"PlayerNum"`
	PlayerName			string				`json:"PlayerName"`
	//MaxHealth			int					`json:"MaxHealth"`
	CurrentHealth		int					`json:"Health"`
	CurrentWeaponId		int					`json:"CurrWeaponId"`
	Inventory			[]*Item				`json:"Inventory,omitempty"`
}