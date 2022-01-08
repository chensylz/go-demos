package main

import (
	"fmt"
	"os/exec"
)

func main() {
	for {
		cmd := exec.Command("ssh", "-N", "-D", "7070", "chenjieping@gate.zaihui.com.cn")

		if err := cmd.Start(); err != nil {
			fmt.Println("Execute failed when Start:" + err.Error())
			return
		}
		fmt.Println("已连接...")
		if err := cmd.Wait(); err != nil {
			fmt.Println("Execute failed when Wait:" + err.Error())
			return
		}
		fmt.Println("断线了,重连中...")
	}
}
