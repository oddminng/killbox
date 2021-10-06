package main

import (
	"flag"
	"fmt"
	"github.com/oddminng/killbox/GameClient"
)

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8000

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8000, "设置服务器端口(默认是8000)")
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

	// 启动客户端的业务
	client.Run()
}
