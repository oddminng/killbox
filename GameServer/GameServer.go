package GameServer

import (
	"fmt"
	"net"
)

type GameServer struct {
	Port int
}

func NewGameServer(port int) *GameServer {
	server := &GameServer{
		Port: port,
	}
	return server
}

func (gameServer *GameServer) Handler(conn net.Conn)  {
	fmt.Println("new conn...")
}

func (gameServer *GameServer) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d",gameServer.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
	}

	defer listener.Close()

	for  {
		conn,err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		go gameServer.Handler(conn)
	}

}
