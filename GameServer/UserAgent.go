package GameServer

import "net"

type UserAgent struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUserAgent(conn net.Conn) *UserAgent {
	userAddr := conn.RemoteAddr().String()
	user := &UserAgent{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	go user.ListenMessage()

	return user
}

func (u UserAgent) ListenMessage() {
	for {
		msg := <- u.C
		u.conn.Write([]byte(msg+"\n"))
	}
}