package model

// 맵 데이터는 좌표 값인 X,Y로 이루어진 수 많읁 좌표값과
// 이미지 랜더링 파일을 위한, 이미지 파일의 정보와 전체 맵 데이터
type MapObj struct {
	tiles	[][]tile
}

// 타일에 대한 정보
type tile struct {
	tileType		int

}