package main

import (
	"fmt"
	"github.com/oddminng/killbox/GameServer"
	"github.com/oddminng/killbox/WebServer"
)

func main() {
	fmt.Println("Start Main")

	go func() {
		fmt.Println("Start WebServer")
		webServer := WebServer.NewWebServer(8080)
		webServer.Start()
	}()
	go func() {
		fmt.Println("Start GameServer")
		gameServer := GameServer.NewGameServer(8000)
		gameServer.Start()
	}()

	//阻塞 main 函数
	select {}
}
