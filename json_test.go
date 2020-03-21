package main

import "testing"

func TestJson(t *testing.T) {
	str := `{"service_name":"login","sign":"15E3DF039F5A02BB1A17316976DE8A51",
	"data":{"union_id":"200128","park_id":"24155",
	"local_id":"c61562ee71b6_2.1.0.0_channels_1_2_3_4_5_6_7_8_9_10_11_12_13_14_15_16",
	"version":"2.1.0.0"}}`
	by := []byte(str)
	aa := parseJSON(by)
	t.Log(aa)
}
