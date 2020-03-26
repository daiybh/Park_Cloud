package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		fmt.Fprintln(c.Writer, "welcome to go website"+time.Now().UTC().String())
		for k, v := range ps.connMap {
			fmt.Fprintln(c.Writer, k, v, v.RemoteAddr().String())
		}
	})
	router.POST("/park/deliverTicket", func(c *gin.Context) {

		contentType := c.ContentType()
		if contentType != "application/json" {
			c.JSON(http.StatusNoContent, gin.H{
				"state":  204,
				"result": "don't support " + contentType,
			})
			return
		}
		fmt.Println("contentType:---->", contentType)
		sign := c.Query("sign")

		fmt.Println("sign->", sign)
		data, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println("ctx.Request.body:", string(data))
		logger.Info("deliverTicket", string(data))
		checkSign(string(data), sign)
		//{"ticket_id":"10022","create_time":"1490879218","money":"5","car_number":"川AD12345","order_id":"9880","remark":"32","park_id":"test001"}
		type DeliverTicket struct {
			TicketID   string `json:"ticket_id"`
			CreateTime string `json:"create_time"`
			Money      string `json:"money"`
			CarNumber  string `json:"car_number"`
			OrderID    string `json:"order_id"`
			Remark     string `json:"remark"`
			ParkID     string `json:"park_id"`
		}
		var ticket DeliverTicket
		json.Unmarshal(data, &ticket)
		fmt.Println("Ticket: ", ticket)

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"age":    123,
			"sign":   sign,
			"token":  ticket.TicketID,
			"state":  200,
			"result": "success",
		})
	})
	//":"+strconv.Itoa(Config.ServerConfig.Httpport)
	//os.Getenv("PORT")
	router.Run(":" + strconv.Itoa(Config.ServerConfig.Httpport))
	logger.Info("End main.......")
}
