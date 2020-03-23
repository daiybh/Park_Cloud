package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type configs struct {
	ServerConfig struct {
		Httpport   int `json:"httpport"`
		Socketport int `json:"socketport"`
	} `json:"serverConfig"`
}

var Config configs

func loadConfig() error {
	filename := "./logs.config"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	//c := configs{}
	err = json.Unmarshal(bytes, &Config)
	if err != nil {
		return err
	}
	return nil
}
func main() {
	loadConfig()
	err := logger.NewLogger("default")
	if err != nil {
		log.Fatal(err)
	}

	defer logger.Close()

	logger.Info("start main.......http:", Config.ServerConfig.Httpport, "  socket:", Config.ServerConfig.Socketport)

	// 创建监听
	ps := ParkServer{}
	ps.HandleFunc("login", func(buf []byte, n int, conn net.Conn) {

	})
	ps.startServer(Config.ServerConfig.Socketport)

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

	err = http.ListenAndServe(":"+string(Config.ServerConfig.Httpport), nil)
	if err != nil {
		logger.Fatal("http.listern ", Config.ServerConfig.Httpport, " failed.", err)
	}
	logger.Info("End main.......")
}
