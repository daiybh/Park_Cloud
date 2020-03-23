package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/kuxuee/logger"
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
	err := logger.NewLogger("default")
	if err != nil {
		log.Fatal(err)
	}

	defer logger.Close()
	// 创建监听
	ps := ParkServer{}
	ps.startServer()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//var s string
		//s = "welcome " + time.Now().UTC().String() + "<br>"
		fmt.Fprintln(w, "welcome to go website"+time.Now().UTC().String())
		for k, v := range ps.connMap {
			fmt.Fprintln(w, k, v, v.RemoteAddr().String())
		}

	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", nil)

}
