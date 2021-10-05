package GameServer

import "net"

type UserAgent struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	gameServer *GameServer
}

func NewUserAgent(conn net.Conn, gameServer *GameServer) *UserAgent {
	userAddr := conn.RemoteAddr().String()
	user := &UserAgent{
		Name:       userAddr,
		Addr:       userAddr,
		C:          make(chan string),
		conn:       conn,
		gameServer: gameServer,
	}

	go user.ListenMessage()

	return user
}

// Online 用户上线
func (u *UserAgent) Online() {
	// 用户上线,创建用户代理,并加入到Server的用户代理Map中
	u.gameServer.mapLock.Lock()
	u.gameServer.UserAgentMap[u.Name] = u
	u.gameServer.mapLock.Unlock()

	// 广播用户上线消息
	u.gameServer.BroadCast(u, "已经上线")
}

// Offline 用户下线
func (u *UserAgent) Offline() {
	// 用户下线，将用户从Server的用户代理Map中删除
	u.gameServer.mapLock.Lock()
	delete(u.gameServer.UserAgentMap, u.Name)
	u.gameServer.mapLock.Unlock()

	// 广播用户下线消息
	u.gameServer.BroadCast(u, "已经下线")
}

// DoMessage 处理用户消息
func (u *UserAgent) DoMessage(msg string) {
	u.gameServer.BroadCast(u, msg)
}

func (u *UserAgent) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
