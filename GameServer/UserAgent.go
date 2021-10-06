package GameServer

import (
	"fmt"
	"net"
	"strings"
)

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

// SendMsg 给当前用户发送消息
func (u *UserAgent) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

// DoMessage 处理用户消息
func (u *UserAgent) DoMessage(msg string) {
	if msg == "who" {
		u.gameServer.mapLock.Lock()
		for _, user := range u.gameServer.UserAgentMap {
			sendMsg := fmt.Sprintf("[%s]%s:在线...\n", user.Addr, user.Name)
			u.SendMsg(sendMsg)
		}
		u.gameServer.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, isOk := u.gameServer.UserAgentMap[newName]
		if isOk {
			u.SendMsg("用户名已被使用\n")
		} else {
			u.gameServer.mapLock.Lock()
			delete(u.gameServer.UserAgentMap, u.Name)
			u.gameServer.UserAgentMap[newName] = u
			u.gameServer.mapLock.Unlock()

			u.Name = newName
			u.SendMsg(fmt.Sprintf("您已经更新用户名：%s\n", u.Name))
		}
	} else if len(msg) > 3 && msg[:3] == "to|" {
		msgSlice := strings.Split(msg, "|")
		remoteName := msgSlice[1]
		if remoteName == "" || len(msgSlice) != 3 {
			u.SendMsg("消息格式不正确，请使用\"to|张三|你好啊\"格式\n")
			return
		}
		msgContent := msgSlice[2]
		if msgContent == "" {
			u.SendMsg("无消息内容，请重发\n")
			return
		}
		remoteUser, isOk := u.gameServer.UserAgentMap[remoteName]
		if !isOk {
			u.SendMsg("该用户不存在\n")
			return
		}
		remoteUser.SendMsg(fmt.Sprintf("[%s]对您说:%s", u.Name, msgContent))
	} else {
		u.gameServer.BroadCast(u, msg)
	}
}

func (u *UserAgent) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
