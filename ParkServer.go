package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/daiybh/logger"
	"github.com/gin-gonic/gin"
)

type ParkServer struct {
	socketServer ParkSocketServer
}

func (ps *ParkServer) run() {
	// 创建监听
	ps.socketServer = ParkSocketServer{}
	ps.socketServer.HandleFunc("login", func(buf []byte, n int, conn net.Conn) {

	})
	ps.socketServer.startServer(Config.ServerConfig.Socketport)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		fmt.Fprintln(c.Writer, "welcome to go website"+time.Now().UTC().String())
		for k, v := range ps.socketServer.connMap {
			fmt.Fprintln(c.Writer, k, v, v.RemoteAddr().String())
		}
	})
	router.GET("/sendTicket/:name", func(c *gin.Context) {
		name := c.Param("name")
		vjson := `{"ticket_type":2,"create_time":1544248573,"limit_day":9999999999,"service_name":"deliver_ticket","have_order":1,
		"park_id":"24155","endtime":"23:59","remark":"","discount":"","starttime":"00:00","ticket_id":"6547339",
		"shop_name":"测全免","startdate":1544248567,"duration":0,"enddate":1563321600,
		"money":"","time_range":0,"car_number":"苏AQW888","order_id":"102"}
		`
		if name != "" {
			vjson = strings.Replace(vjson, "苏AQW888", name, 1)
		}
		b := ClientGroup.SendToClient("24155", []byte(vjson))

		c.JSON(http.StatusOK, gin.H{
			"state":  200,
			"result": b,
			"vv":     vjson,
		})
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
		sign := c.Query("sign")
		data, _ := ioutil.ReadAll(c.Request.Body)
		logger.Info("deliverTicket", string(data))
		if !checkSign(string(data), sign) {
			c.JSON(http.StatusOK, gin.H{
				"state":  -1001007,
				"result": "wrong sign",
			})
			return
		}
		//{"ticket_id":"10022","create_time":"1490879218","money":"5","car_number":"川AD12345","order_id":"9880","remark":"32","park_id":"test001"}

		var ticket JumaDeliverTicket
		json.Unmarshal(data, &ticket)
		HandleJuMaTick(ticket)
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"age":    123,
			"sign":   sign,
			"token":  ticket.TicketID,
			"state":  200,
			"result": "success",
		})
	})
	router.Run(":" + strconv.Itoa(Config.ServerConfig.Httpport))
	logger.Info("End main.......")
}
