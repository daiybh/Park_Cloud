package main

import (
	"fmt"
	"net"
	"strings"
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
func HandleConn(conn net.Conn) {
	defer conn.Close()
	//获取客户端的网络地址信息
	addr := conn.RemoteAddr().String()
	fmt.Println(addr, " conncet sucessful")

	buf := make([]byte, 2048)

	for {
		//读取用户数据
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		fmt.Println(string(buf[:n]))
		ret, er := parseJSON(buf[:n])
		if er == nil {
			fmt.Println("goback-->", ret)
			time.Sleep(1 * time.Second)
			n, err = conn.Write([]byte(ret + "\r\n"))
			fmt.Println("write", n, err)
			if err != nil {
				return
			}
			//go heartPack(conn)
		} else {
			//把数据转换为大写，再给用户发送
			fmt.Println("nos")
			conn.Write([]byte(strings.ToUpper(string(buf[:n]))))
		}
	}

}
func main() {

	//fmt.Println(m["servers"])
	// 创建监听
	listener, err := net.Listen("tcp", ":6789")
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	defer listener.Close() // 主协程结束时，关闭listener

	for {
		fmt.Println("服务器等待客户端建立连接...")

		// 等待客户端连接请求
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
			return
		}
		fmt.Println("客户端与服务器连接建立成功...")
		go HandleConn(conn)
	}

}
