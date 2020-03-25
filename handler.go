package main

import (
	"encoding/json"
	"net"
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
	conn   net.Conn
	parkID string
	token  string
}
type _ClientGroup struct {
	m  map[string]_ClientInfo //key:park_id value: Net.conn
	mu sync.Mutex
}

func (g *_ClientGroup) GetToken(parkID string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	v, ok := g.m[parkID]
	if ok {
		return v.token
	}

	if g.m == nil {
		g.m = make(map[string]_ClientInfo)
	}
	e := _ClientInfo{parkID: parkID, token: strings.ToLower(MD5(parkID))}
	g.m[parkID] = e
	return e.token
}
func (g *_ClientGroup) CheckToken(parkID string, token string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	v, ok := g.m[parkID]
	if ok {
		return 0 == strings.Compare(strings.ToLower(token), v.token)
	}
	return false
}

var ClientGroup = _ClientGroup{}

func NeedRecord(carNumber string) bool {
	logger.Info("carNumber:", carNumber)
	return true
}
func handleLogin(buf []byte, n int) string {

	parkID := jsoniter.Get(buf, "data", "park_id").ToString()
	token := ClientGroup.GetToken(parkID)

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

	if err := jsoniter.Unmarshal(jsonStr, &inPark); err != nil {
		logger.Error("Unmarshal failed.", jsonStr)
		return ""
	}

	if !ClientGroup.CheckToken(inPark.Data.ParkID, inPark.token) {
		logger.Error("wrong token ", inPark.Data.ParkID, inPark.token)
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
	if err := jsoniter.Unmarshal(jsonStr, &outPark); err != nil {
		logger.Error("Unmarshal failed.", jsonStr)
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
