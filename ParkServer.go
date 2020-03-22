package main

import (
	"fmt"
	"net"
	"time"
)

type ParkServer struct {
	listener net.Listener
	connMap  map[string]net.Conn //key:park_id value: Net.conn
}

func (ps *ParkServer) startServer() {
	fmt.Println("package park")
	var err error
	ps.listener, err = net.Listen("tcp", ":6789")
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	defer ps.listener.Close() // 主协程结束时，关闭listener
	ps.connMap = make(map[string]net.Conn)
	for {
		fmt.Println("服务器等待客户端建立连接...")

		// 等待客户端连接请求
		conn, err := ps.listener.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
			return
		}
		fmt.Println("客户端与服务器连接建立成功...")

		go ps.HandleConn(conn)
	}
}

func (ps *ParkServer) HandleConn(conn net.Conn) {
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
		er, serverName, parkid, ret := parseJSON(buf[:n])
		if er == nil {
			fmt.Println("goback-->", serverName, parkid, ret)

			if serverName == "login" {
				ps.connMap[parkid] = conn
			}
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
			//	conn.Write([]byte(strings.ToUpper(string(buf[:n]))))
		}
	}

}
