package main

import (
	"encoding/json"
	"net"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

type loginReturn struct {
	State       int    `json:"state"`
	Token       string `json:"token"`
	ServiceName string `json:"service_name"`
}
type inParkReturn struct {
	ServiceName string `json:"service_name"`
	ParkID      string `json:"park_id"`
	Errmsg      string `json:"errmsg"`
	State       int    `json:"state"`
	OrderID     string `json:"order_id"`
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
	connMap map[string]net.Conn //key:park_id value: Net.conn
	mu      sync.Mutex
}

var clientInfos = &_ClientInfo{}

func handleLogin(buf []byte, n int) string {
	lret := loginReturn{
		ServiceName: "login",
		State:       1,
		Token:       "5880277f494544259642dd7ac35afdf4",
	}
	b, _ := json.Marshal(lret)

	return string(b) // `{"state":1,"token":"5880277f494544259642dd7ac35afdf4","service_name":"login"}`
}

func handleInPark(jsonStr []byte, n int) string {
	parkID := jsoniter.Get(jsonStr, "data", "park_id").ToString()
	orderID := jsoniter.Get(jsonStr, "data", "order_id").ToString()
	in := inParkReturn{
		ParkID:      parkID,
		Errmsg:      "",
		State:       1,
		OrderID:     orderID,
		ServiceName: "in_park",
	}
	b, _ := json.Marshal(in)

	return string(b) // `{"state":1,"token":"5880277f494544259642dd7ac35afdf4","service_name":"login"}`
}

func handleOutpark(jsonStr []byte, n int) string {
	//parkID := jsoniter.Get(jsonStr, "data", "park_id").ToString()
	orderID := jsoniter.Get(jsonStr, "data", "order_id").ToString()
	payType := jsoniter.Get(jsonStr, "data", "pa_type").ToString()
	outRet := outParkReturn{
		Errmsg:      "",
		State:       1,
		OrderID:     orderID,
		ServiceName: "out_park",
		PayType:     payType,
	}
	b, _ := json.Marshal(outRet)

	return string(b) // `{"state":1,"token":"5880277f494544259642dd7ac35afdf4","service_name":"login"}`
}
