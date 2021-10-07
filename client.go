package main

import (
	"flag"
	"fmt"
	"github.com/oddminng/killbox/GameClient"
)

var serverIp string
var serverPort int
var userName string
var userPass string

// ./client -ip 127.0.0.1 -port 8000

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8000, "设置服务器端口(默认是8000)")
	flag.StringVar(&userName, "user", "", "设置登录用户名(10001~99999)")
	flag.StringVar(&userPass, "pw", "123", "设置登录密码(默认是123)")
}

func main() {
	// 命令行解析
	flag.Parse()

	fmt.Println(">>>>>> Start Client")
	client := GameClient.NewGameClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>> 连接服务器失败...")
	}

	// 开启处理服务端回执的 goroutine
	go client.DealResponse()
	fmt.Println(">>>>>> 连接服务端成功...")

	if userName != "" && userPass != "" {
		client.Login(userName, userPass)
		fmt.Println(">>>>>> 登录...")
	}

	// 启动心跳包发送
	client.KeepAlive()

	// 启动客户端的业务
	client.Run()
}
