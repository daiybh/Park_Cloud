package main

import (
	"testing"
)

func TestJsonlogin(t *testing.T) {
	str := `{"service_name":"login","sign":"15E3DF039F5A02BB1A17316976DE8A51",
	"data":{"union_id":"200128","park_id":"24155",
	"local_id":"c61562ee71b6_2.1.0.0_channels_1_2_3_4_5_6_7_8_9_10_11_12_13_14_15_16",
	"version":"2.1.0.0"}}`
	by := []byte(str)
	aa := handleLogin(by, len(str))
	t.Log(" aa:", aa)

}

func TestJson_inpark(t *testing.T) {
	str := `{"service_name":"in_park","sign":"987B2045CDCFF2FAFDA392E3EA8093B4","token":"5880277f494544259642dd7ac35afdf4",
	"data":{"car_number":"绮W4444","in_time":1577244491,"order_id":"302","empty_plot":885,"park_id":"24155"}}
	`
	by := []byte(str)
	aa := handleInPark(by, len(str))

	t.Log(" aa:", aa)
}

func TestJson_outpark(t *testing.T) {
	str := `{"service_name":"out_park","sign":"DD0BD8EAFE672B4741B4F3F523E794F3","token":"5880277f494544259642dd7ac35afdf4","data":{"car_number":"粤B1H7S0","in_time":1576327327,"out_time":1576327362,"total":"0.0","order_id":"1131522704","empty_ plot":0,"park_id":"24155","pay_type":"cash","auth_code":""}}
	`
	by := []byte(str)
	aa := handleOutpark(by, len(str))

	t.Log(" aa:", aa)
}

func TestJson_XXXXX(t *testing.T) {
	str := `{"service_namea":"out_park","sign":"DD0BD8EAFE672B4741B4F3F523E794F3","token":"5880277f494544259642dd7ac35afdf4","data":{"car_number":"粤B1H7S0","in_time":1576327327,"out_time":1576327362,"total":"0.0","order_id":"1131522704","empty_ plot":0,"park_id":"24155","pay_type":"cash","auth_code":""}}
	`
	by := []byte(str)
	aa := handleInPark(by, len(str))

	t.Log(" aa:", aa)
}

func TestCallBack(t *testing.T) {
	JumaTest()
	t.Log(" aa:", "aa")
}
