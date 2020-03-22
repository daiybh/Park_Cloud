package main

import (
	"fmt"
	"net"
	"time"
)

func heartPack(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1)
	//buf := []byte(0x11)
	buf[0] = 0x11
	count := 0
	for {
		n, err := conn.Write(buf)
		if err != nil {
			fmt.Println("hearPack err = ", err)
			break
		}
		fmt.Println(time.Now().UTC(), "send heartPack", count, n, err)
		count++
		time.Sleep(1 * time.Second)
	}
}

func main() {
	//fmt.Println(m["servers"])
	// 创建监听
	ps := ParkServer{}
	ps.startServer()

}
