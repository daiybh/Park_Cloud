package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/kuxuee/logger"
)

type loginReturn struct {
	State       int    `json:"state"`
	Token       string `json:"token"`
	ServiceName string `json:"service_name"`
}
type inParkCome struct {
	ServiceName string `json:"service_name"`
	Sign        string `json:"sign"`
	Token       string `json:"token"`
	Data        struct {
		CarNumber string `json:"car_number"`
		InTime    int    `json:"in_time"`
		OrderID   string `json:"order_id"`
		EmptyPlot int    `json:"empty_plot"`
		ParkID    string `json:"park_id"`
	} `json:"data"`
}

type inParkReturn struct {
	ServiceName string `json:"service_name"`
	ParkID      string `json:"park_id"`
	Errmsg      string `json:"errmsg"`
	State       int    `json:"state"`
	OrderID     string `json:"order_id"`
}
type outParkCome struct {
	ServiceName string `json:"service_name"`
	Sign        string `json:"sign"`
	Token       string `json:"token"`
	Data        struct {
		CarNumber string `json:"car_number"`
		InTime    int    `json:"in_time"`
		OutTime   int    `json:"out_time"`
		Total     string `json:"total"`
		OrderID   string `json:"order_id"`
		EmptyPlot int    `json:"empty_ plot"`
		ParkID    string `json:"park_id"`
		PayType   string `json:"pay_type"`
		AuthCode  string `json:"auth_code"`
	} `json:"data"`
}
type outParkReturn struct {
	ServiceName string `json:"service_name"`
	OrderID     string `json:"order_id"`
	PayType     string `json:"pay_type"`
	NetStatus   int    `json:"net_status"`
	State       int    `json:"state"`
	Errmsg      string `json:"errmsg"`
}
type _ClientInfo struct {
	Conn   net.Conn
	ParkID string
	Token  string
}
type _ClientGroup struct {
	m  map[string]_ClientInfo //key:park_id value: Net.conn
	mu sync.Mutex
}

func (g *_ClientGroup) Login(conn net.Conn, parkID string) string {

	g.mu.Lock()
	defer g.mu.Unlock()
	v, ok := g.m[parkID]
	if ok {
		v.Conn = conn
		return v.Token
	}

	if g.m == nil {
		g.m = make(map[string]_ClientInfo)
	}
	e := _ClientInfo{ParkID: parkID, Conn: conn, Token: strings.ToLower(MD5(parkID))}
	g.m[parkID] = e
	return e.Token
}
func (g *_ClientGroup) CheckToken(parkID string, token string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	v, ok := g.m[parkID]
	if ok {
		return 0 == strings.Compare(strings.ToLower(token), v.Token)
	}
	return false
}
func (g *_ClientGroup) SendToClient(parkID string, msg []byte) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	v, ok := g.m[parkID]
	if !ok || v.Conn == nil {
		logger.Error("cannot found conn ", parkID, msg)
		return false
	}
	_, err := v.Conn.Write(msg)
	if err != nil {
		v.Conn = nil
	}
	logger.Debug("sendToClient-->", string(msg))
	//v.Conn.Write([]byte("\r\n"))
	return err == nil
}

var ClientGroup = _ClientGroup{}

func NeedRecord(carNumber string) bool {
	logger.Info("carNumber:", carNumber)
	return true
}
func handleLogin(conn net.Conn, buf []byte, n int) string {

	parkID := jsoniter.Get(buf[:n], "data", "park_id").ToString()
	token := ClientGroup.Login(conn, parkID)

	lret := loginReturn{
		ServiceName: "login",
		State:       1,
		Token:       token,
	}
	b, _ := json.Marshal(lret)

	return string(b) // `{"state":1,"token":"5880277f494544259642dd7ac35afdf4","service_name":"login"}`
}

func handleInPark(jsonStr []byte, n int) string {
	//str := `{"service_name":"in_park","sign":"987B2045CDCFF2FAFDA392E3EA8093B4","token":"5880277f494544259642dd7ac35afdf4","data":{"car_number":"绮W4444","in_time":1577244491,"order_id":"302","empty_plot":885,"park_id":"24155"}}`

	var inPark inParkCome

	if err := jsoniter.Unmarshal(jsonStr[:n], &inPark); err != nil {
		logger.Error("Unmarshal failed.", string(jsonStr[:n]))
		return ""
	}

	if !ClientGroup.CheckToken(inPark.Data.ParkID, inPark.Token) {
		logger.Error("wrong token ", inPark.Data.ParkID, inPark.Token)
		return ""
	}

	if NeedRecord(inPark.Data.CarNumber) {

	}

	JumaMakeParkIn()
	in := inParkReturn{
		ParkID:      inPark.Data.ParkID,
		Errmsg:      "",
		State:       1,
		OrderID:     inPark.Data.OrderID,
		ServiceName: "in_park",
	}
	b, _ := json.Marshal(in)

	return string(b)
}

func handleOutpark(jsonStr []byte, n int) string {
	//	str := `{"service_name":"out_park","sign":"DD0BD8EAFE672B4741B4F3F523E794F3","token":"5880277f494544259642dd7ac35afdf4","data":{"car_number":"粤B1H7S0","in_time":1576327327,"out_time":1576327362,"total":"0.0","order_id":"1131522704","empty_ plot":0,"park_id":"24155","pay_type":"cash","auth_code":""}}`

	var outPark outParkCome
	if err := jsoniter.Unmarshal(jsonStr[:n], &outPark); err != nil {
		logger.Error("Unmarshal failed.", jsonStr[:n])
		return ""
	}
	if !ClientGroup.CheckToken(outPark.Data.ParkID, outPark.Token) {
		logger.Error("wrong token ", outPark.Data.ParkID, outPark.Token)
		return ""
	}
	if NeedRecord(outPark.Data.CarNumber) {

	}
	JumaMakeParkOut()
	outRet := outParkReturn{
		Errmsg:      "",
		State:       1,
		OrderID:     outPark.Data.OrderID,
		ServiceName: "out_park",
		PayType:     outPark.Data.PayType,
	}
	b, _ := json.Marshal(outRet)

	return string(b)
}

type Ticket struct {
	TicketType  int    `json:"ticket_type"`
	CreateTime  int    `json:"create_time"`
	LimitDay    int64  `json:"limit_day"`
	ServiceName string `json:"service_name"`
	HaveOrder   int    `json:"have_order"`
	ParkID      string `json:"park_id"`
	Endtime     string `json:"endtime"`
	Remark      string `json:"remark"`
	Discount    string `json:"discount"`
	Starttime   string `json:"starttime"`
	TicketID    string `json:"ticket_id"`
	ShopName    string `json:"shop_name"`
	Startdate   int    `json:"startdate"`
	Duration    int    `json:"duration"`
	Enddate     int    `json:"enddate"`
	Money       string `json:"money"`
	TimeRange   int    `json:"time_range"`
	CarNumber   string `json:"car_number"`
	OrderID     string `json:"order_id"`
}

func HandleJuMaTick(jumaTick JumaDeliverTicket) {
	//find the client from clientGroups
	//transform struct
	//send
	createTime, _ := strconv.Atoi(jumaTick.CreateTime)

	tick := Ticket{
		OrderID:     jumaTick.OrderID,
		TicketID:    jumaTick.TicketID,
		Money:       jumaTick.Money,
		ParkID:      jumaTick.ParkID,
		CarNumber:   jumaTick.CarNumber,
		Remark:      jumaTick.Remark,
		CreateTime:  createTime,
		Starttime:   "00:00",
		Endtime:     "23:59",
		Startdate:   createTime,
		Enddate:     createTime + 100,
		Duration:    0,
		TimeRange:   0,
		HaveOrder:   1,
		LimitDay:    99999999,
		ServiceName: "deliver_ticket",
	}
	vbyte, _ := json.Marshal(tick)

	vs := string(vbyte)
	vs += "\r\n"
	ClientGroup.SendToClient(tick.ParkID, []byte(vs))
	fmt.Println(tick)
}
