package main

import (
	"fmt"
	"net"

	jsoniter "github.com/json-iterator/go"
	"github.com/kuxuee/logger"
)

type ParkServer struct {
	listener net.Listener
	connMap  map[string]net.Conn //key:park_id value: Net.conn
	m        map[string]muxEntry
}
type Handler func(buf []byte, n int, conn net.Conn) string
type muxEntry struct {
	h       Handler
	pattern string
}
type HandlerFunc func([]byte, int, net.Conn)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(buf []byte, n int, conn net.Conn) {
	f(buf, n, conn)
}
func (ps *ParkServer) HandleFunc(pattern string, handler func(buf []byte, n int, conn net.Conn)) {
	//	DefaultServeMux.HandleFunc(pattern, handler)
	if ps.m == nil {
		ps.m = make(map[string]muxEntry)
	}
	//hh := HandlerFunc(handler)

	//e := muxEntry{h: hh, pattern: pattern}
	//ps.m[pattern] = e
}
func (ps *ParkServer) acceptThread() {
	defer ps.listener.Close() // 主协程结束时，关闭listener
	fmt.Println("parkServer   acceptThrad....")
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
func (ps *ParkServer) startServer() {
	fmt.Println("package park")
	var err error
	ps.listener, err = net.Listen("tcp", ":6789")
	if err != nil {
		logger.Fatal("listen err:", err)
		return
	}
	logger.Info("startServer....")
	go ps.acceptThread()
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
		ps.handleMessage(conn, buf, n)

	}

}

func (ps *ParkServer) handleMessage(conn net.Conn, buf []byte, n int) {

	servicename := jsoniter.Get(buf, "service_name").ToString()
	fmt.Println("servicename ---->", servicename)
	logger.Info(string(buf[:n]))
	if servicename == "" {
		fmt.Println("servicename null")
		return
	}
	token := jsoniter.Get(buf, "token").ToString()
	fmt.Println("handleMessage", servicename, token)

	var retString string
	switch servicename {
	case "login":
		retString = handleLogin(buf, n)
	case "in_park":
		retString = handleInPark(buf, n)
	case "out_park":
		retString = handleOutpark(buf, n)

	}
	if retString == "" {
		return
	}
	n, err := conn.Write([]byte(retString + "\r\n"))
	if err != nil {
		return
	}

}
