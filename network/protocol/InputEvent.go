package protocol

type InputData struct {
	PlayerNum		int			`json:"PlayerNum"`
	AxisX			int			`json:"AxisX"`
	AxisY			int			`json:"AxisY"`
	PositionX		float32		`json:"PositionX"`
	PositionY		float32		`json:"PositionY"`
	MouseX			int			`json:"MouseX"`
	MouseY			int			`json:"MouseY"`
}