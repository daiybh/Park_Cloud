package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func checkSign(jsonData string, sign string) bool {
	testKey := "206d98f4de9e4423b3aa42cd0c36fd84"
	ss := jsonData + testKey
	md5String := MD5(ss)
	x := strings.Compare(md5String, sign)
	fmt.Println("comapre-->", x, md5String, sign)
	return x == 0
}
func makepostURL(actionName string, sValue string) string {
	testKey := "206d98f4de9e4423b3aa42cd0c36fd84"
	ss := sValue[:len(sValue)-1] + testKey
	md5String := MD5(ss)
	baseURL := "http://test.api.cxhshop.xyz/parkmall/" + actionName
	return baseURL + "?sign=" + md5String
}

//{"ParkCode":"5101070030","VehicleNo":"川DDDD77","StartTime":"2019-08-09 15:36:21"}
type JumaParkIn struct {
	ParkCode  string `json:"ParkCode"`
	VehicleNo string `json:"VehicleNo"`
	StartTime string `json:"StartTime"`
}

//{"ParkCode":"5101070030","VehicleNo":"川DDDD77","StartTime":"2019-08-09 15:36:21","EndTime":"","PaymentMoney":1234.23}
type JumaParkOut struct {
	ParkCode     string  `json:"ParkCode"`
	VehicleNo    string  `json:"VehicleNo"`
	StartTime    string  `json:"StartTime"`
	EndTime      string  `json:"EndTime"`
	PaymentMoney float64 `json:"PaymentMoney"`
}
type JumaDeliverTicket struct {
	TicketID   string `json:"ticket_id"`
	CreateTime string `json:"create_time"`
	Money      string `json:"money"`
	CarNumber  string `json:"car_number"`
	OrderID    string `json:"order_id"`
	Remark     string `json:"remark"`
	ParkID     string `json:"park_id"`
}

func parkPost(actionName string, b *bytes.Buffer) {
	newURL := makepostURL(actionName, b.String())
	fmt.Println(newURL)
	resp, err := http.Post(newURL, "application/json", b)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	fmt.Println(string(body))
}
func JumaMakeParkIn() {
	u := JumaParkIn{ParkCode: "5101070030", VehicleNo: "川DDDD77", StartTime: "2019-08-09 15:36:21"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	parkPost("carStart", b)
}
func JumaMakeParkOut() {
	u := JumaParkOut{ParkCode: "5101070030", VehicleNo: "川DDDD77", StartTime: "2019-08-09 15:36:21", EndTime: "2019-08-09 19:36:21", PaymentMoney: 19.8}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	parkPost("carEnd", b)
}
func JumaTest() {
	JumaMakeParkIn()
	JumaMakeParkOut()
}
