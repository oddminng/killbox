package GameClient

import (
	"fmt"
	"io"
	"net"
	"os"
)

type GameClient struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // 当前客户端模式
}

func NewGameClient(serverIp string, serverPort int) *GameClient {
	// 创建客户端对象
	client := &GameClient{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	// 连接Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("new.Dial error", err)
		return nil
	}
	client.conn = conn

	// 返回对象
	return client
}

func (c *GameClient) DealResponse() {
	if _, err := io.Copy(os.Stdout, c.conn); err != nil {
		fmt.Println("接收服务端消息失败 err:", err)
		return
	}
}

func (c *GameClient) menu() bool {
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	_, err := fmt.Scanln(&c.flag)
	if err != nil {
		fmt.Println(">>>> 解析输入失败 err:", err)
		return false
	}

	if c.flag >= 0 && c.flag <= 3 {
		return true
	} else {
		fmt.Println(">>>> 请输入合法范围内的数字 <<<<")
		return false
	}
}

func (c *GameClient) updateName() bool {
	fmt.Println(">>>>请输入用户名:")
	if _, err := fmt.Scanln(&c.Name); err != nil {
		fmt.Println(">>>> 解析输入失败 err:", err)
		return false
	}

	sendMsg := fmt.Sprintf("rename|%s\n", c.Name)

	if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

func (c *GameClient) publicChat() {
	// 提示用户输入消息
	var chatMsg string

	fmt.Println(">>>> 请输入聊天内容, exit 退出.")
	if _, err := fmt.Scanln(&chatMsg); err != nil {
		fmt.Println("解析输入失败 err:", err)
		return
	}

	for chatMsg != "exit" {

		// 发送给服务器
		if len(chatMsg) != 0 {
			sendMsg := fmt.Sprintf("%s\n", chatMsg)
			if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>> 请输入聊天内容, exit 退出.")
		if _, err := fmt.Scanln(&chatMsg); err != nil {
			fmt.Println("解析输入失败 err:", err)
			return
		}
	}

}

func (c *GameClient) showUsers() {
	sendMsg := "who\n"
	if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
		fmt.Println("conn Write err:", err)
	}
}

func (c *GameClient) privateChat() {
	var remoteName string
	var chatMsg string

	c.showUsers()
	fmt.Println(">>>> 请输入聊天对象的[用户名], exit 退出:")

	if _, err := fmt.Scanln(&remoteName); err != nil {
		fmt.Println("解析输入失败 err:", err)
		return
	}

	if remoteName != "exit" {
		fmt.Println(">>>> 请入消息内容, exit 退出:")
		if _, err := fmt.Scanln(&chatMsg); err != nil {
			fmt.Println("解析输入失败 err:", err)
			return
		}
		for chatMsg != "exit" {
			// 发送给服务器
			if len(chatMsg) != 0 {
				sendMsg := fmt.Sprintf("to|%s|%s\n\n", remoteName, chatMsg)
				if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>> 请入消息内容, exit 退出:")
			if _, err := fmt.Scanln(&chatMsg); err != nil {
				fmt.Println("解析输入失败 err:", err)
				return
			}
		}

		c.showUsers()
		fmt.Println(">>>> 请输入聊天对象的[用户名], exit 退出:")

		if _, err := fmt.Scanln(&remoteName); err != nil {
			fmt.Println("解析输入失败 err:", err)
			return
		}
	}
}

func (c *GameClient) Run() {
	for c.flag != 0 {
		for c.menu() != true {
		}

		switch c.flag {
		case 1:
			// 公聊模式
			c.publicChat()
			break
		case 2:
			// 私聊模式
			c.privateChat()
			break
		case 3:
			// 更新用户名
			c.updateName()
			break
		}
	}
}
