package GameServer

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type UserAgent struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
	gameServer *GameServer
	isLogin bool
}

func NewUserAgent(conn net.Conn, gameServer *GameServer) *UserAgent {
	userAddr := conn.RemoteAddr().String()
	user := &UserAgent{
		Name:       userAddr,
		Addr:       userAddr,
		C:          make(chan string),
		conn:       conn,
		gameServer: gameServer,
		isLogin: false,
	}

	go user.ListenMessage()

	return user
}

// Online 用户上线
func (u *UserAgent) Online() {
	// 用户上线,创建用户代理,并加入到Server的用户代理Map中
	u.gameServer.mapLock.Lock()
	delete(u.gameServer.UserAgentMap, u.Name)
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
	if u.isLogin {
		// 已登录处理用户业务
		if msg == "keepalive" {
			// 维持活跃状态
			fmt.Println(fmt.Sprintf(">>>> user[%s] keepalive %s", u.Name, time.Now().String()))
			return
		} else if msg == "who" {
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
			if len(msgSlice) != 3 {
				u.SendMsg("消息格式不正确，请使用\"to|张三|你好啊\"格式\n")
				return
			}
			remoteName := msgSlice[1]
			msgContent := msgSlice[2]
			if remoteName == "" {
				u.SendMsg("消息格式不正确，请使用\"to|张三|你好啊\"格式\n")
				return
			}
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
	} else {
		// 未登录只处理登录业务
		if len(msg) > 6 && msg[:6] == "login|" {
			msgSlice := strings.Split(msg, "|")
			if len(msgSlice) != 3 {
				u.SendMsg("消息格式不正确，请使用\"login|username|password\"格式\n")
				return
			}
			userName := msgSlice[1]
			userPassword := msgSlice[2]

			if userPassword != "123" {
				u.SendMsg("密码错误\n")
				return
			}

			// 设置用户信息
			u.Name = userName
			u.isLogin = true

			u.Online()
		}
	}
}

func (u *UserAgent) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
