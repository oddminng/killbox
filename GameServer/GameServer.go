package GameServer

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type GameServer struct {
	// 游戏服务端监听端口
	Port int

	// 在线用户列表
	UserAgentMap map[string]*UserAgent
	mapLock      sync.RWMutex

	// 消息广播 channel
	Message chan string
}

func NewGameServer(port int) *GameServer {
	server := &GameServer{
		Port:         port,
		UserAgentMap: make(map[string]*UserAgent),
		Message:      make(chan string),
	}
	return server
}

// ListenMessage 监听 Message Channel，有消息就发送给全部在线User
func (g *GameServer) ListenMessage() {
	for {
		msg := <-g.Message

		// 将消息发送给全部在线用户
		g.mapLock.Lock()
		for _, user := range g.UserAgentMap {
			user.C <- msg
		}
		g.mapLock.Unlock()
	}
}

func (g *GameServer) BroadCast(userAgent *UserAgent, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s", userAgent.Addr, userAgent.Name, msg)

	g.Message <- sendMsg
}

func (g *GameServer) Handler(conn net.Conn) {
	// 当前连接的业务处理
	fmt.Println("new conn...")

	user := NewUserAgent(conn)

	// 用户上线,创建用户代理,并加入到Server的用户代理Map中
	g.mapLock.Lock()
	g.UserAgentMap[user.Name] = user
	g.mapLock.Unlock()

	// 广播用户上线消息
	g.BroadCast(user, "已经上线")

	// 接收用户消息
	go func() {
		buff := make([]byte, 2048)
		for {
			n, err := conn.Read(buff)
			if n == 0 {
				g.BroadCast(user, "下线")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// 提取消息 去除消息尾部'\n'
			msg := string(buff[:n-1])

			// 广播消息
			g.BroadCast(user, msg)
		}
	}()

	// 阻塞当前 Handler
	select {}
}

func (g *GameServer) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
	}

	// close listen socket
	defer listener.Close()

	// 启动监听Message Channel 的 goroutine
	go g.ListenMessage()

	fmt.Println("Game Server Listening and serving TCP on:", g.Port)

	for {
		// accept 处理新连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// do handler
		go g.Handler(conn)
	}

}