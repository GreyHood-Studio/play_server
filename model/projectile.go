package model

// 발사체 (non targeting 되어지는 총알에 대한 처리)를 위한 발사체 정보
type projectile struct {
	shotter 		int	// 총알을 쏜 사람에 대한 구분자
	startPos		Pos	// 총알 처음 발사 시간에 대한 정보
	ttl 			int	// 총알이 사라지는 시간에 대한 정보
	p_size		int	// 총알 범위 판정을 위해 총알의 범위
	p_prop_type 	int	// 총알에 성질 대한 타입 ( 데미지, 속성(독, 마비) etc... )
	p_move_type	int // 총알의 발사 궤적에 대한 타입 ( 방향, 속도 )
}