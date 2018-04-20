# play_server

유저 로그인 시 (유저 id와 pwd를 받고, validation 체크 필요)

호출 순 ( 순환 참조 방지 )
User -> TCP Network -> Packet -> Controller
api -> Gin Router -> route -> Controller

host/port/serverId/ClientId로 Client들 식별 가능

웹서버
로직 처리 모두 route에 존재

게임 서버
호직은 controller / object에 대한 명세는 model