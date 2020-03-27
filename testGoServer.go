package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/daiybh/logger"
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

	logger.Info("##############startMain###############http:", Config.ServerConfig.Httpport, "  socket:", Config.ServerConfig.Socketport)
	ps := ParkServer{}
	ps.run()
}
