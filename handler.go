package main

import (
	"encoding/json"
	"fmt"
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
func parseJSON(jsonStr []byte) (err error, ret string, actionName string, park_id string) {
	fmt.Println(string(jsonStr))
	servicename := jsoniter.Get(jsonStr, "service_name").ToString()
	parkID := jsoniter.Get(jsonStr, "data", "park_id").ToString()
	token := jsoniter.Get(jsonStr, "token").ToString()
	fmt.Println("parseJSON", servicename, parkID, token)
	var retNil error
	var retString string
	retString = "Error"
	retNil = nil

	if servicename == "login" {
		//{"service_name":"login","sign":"15E3DF039F5A02BB1A17316976DE8A51","data":{"union_id":"200128","park_id":"24155","local_id":"c61562ee71b6_2.1.0.0_channels_1_2_3_4_5_6_7_8_9_10_11_12_13_14_15_16","version":"2.1.0.0"}}
		lret := loginReturn{
			ServiceName: servicename,
			State:       1,
			Token:       "5880277f494544259642dd7ac35afdf4",
		}
		b, _ := json.Marshal(lret)

		retString = string(b) // `{"state":1,"token":"5880277f494544259642dd7ac35afdf4","service_name":"login"}`
	} else if servicename == "in_park" {
		//{"service_name":"in_park","sign":"987B2045CDCFF2FAFDA392E3EA8093B4","token":"5880277f494544259642dd7ac35afdf4","data":{"car_number":"绮W4444","in_time":1577244491,"order_id":"302","empty_plot":885,"park_id":"24155"}}

		orderID := jsoniter.Get(jsonStr, "data", "order_id").ToString()
		in := inParkReturn{
			ParkID:      parkID,
			Errmsg:      "",
			State:       1,
			OrderID:     orderID,
			ServiceName: servicename,
		}
		b, _ := json.Marshal(in)
		//fmt.Println("json.marshal-->", ero, string(b))
		retString = string(b) //`{"service_name":"in_park","park_id":"21807","errmsg":"","state":1,"order_id":"102"}`
	} else if servicename == "out_park" {
		//{"service_name":"out_park","sign":"DD0BD8EAFE672B4741B4F3F523E794F3","token":"5880277f494544259642dd7ac35afdf4","data":{"car_number":"粤B1H7S0","in_time":1576327327,"out_time":1576327362,"total":"0.0","order_id":"1131522704","empty_ plot":0,"park_id":"24155","pay_type":"cash","auth_code":""}}
		orderID := jsoniter.Get(jsonStr, "data", "order_id").ToString()
		payType := jsoniter.Get(jsonStr, "data", "pa_type").ToString()
		outRet := outParkReturn{
			Errmsg:      "",
			State:       1,
			OrderID:     orderID,
			ServiceName: servicename,
			PayType:     payType,
		}
		b, _ := json.Marshal(outRet)
		//fmt.Println("json.marshal-->", ero, string(b))
		retString = string(b) // `{"service_name":"out_park","order_id":"102","pay_type":"cash","net_status":1,"state":1,"errmsg":"操作成功完成。"}`
	} else {
		fmt.Println(servicename, parkID, token)
		retNil = nil //
	}

	return retNil, servicename, parkID, retString
}
func TwoSum(nums []int, target int) []int {
	fmt.Println("nums", nums)
	m := make(map[int]int)
	for i, num := range nums {
		key := target - num
		if j, ok := m[key]; ok {
			return []int{j, i}
		}
		m[nums[i]] = i
	}
	return []int{}
}
