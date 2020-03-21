package main

import (
	"fmt"

	"github.com/buger/jsonparser"
)

func mainjsonparser() {
	str := `{"service_name":"login","sign":"304835709614D645AC943F7AC8EA7762",
   "data":{"union_id":"200128","park_id":"24155",
   "local_id":"421561813b17_2.1.0.0_channels_1_2_3_4_5_6_7_8_9_10_11_12_13_14_15_16","version":"2.1.0.0"}}`

	dbyte := []byte(str)
	re, er := jsonparser.GetString(dbyte, "service_name")
	fmt.Println(er, re)
	data := []byte(`{
		"person": {
		  "name":{
			"first": "Leonid",
			"last": "Bugaev",
			"fullName": "Leonid Bugaev"
		  },
		  "github": {
			"handle": "buger",
			"followers": 109
		  },
		  "avatars": [
			{ "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" },
			
			{ "url": "https://avatars1.baidu.com", "type": "thumbnail" }
		  ]
		},
		"company": {
		  "name": "Acme"
		}
	  }`)

	result, err := jsonparser.GetString(data, "person", "name", "fullName")
	if err != nil {
		fmt.Println("err1:", err)
	}
	fmt.Println(result)
	content, valueType, offset, err := jsonparser.Get(data, "person", "name", "fullName")
	if err != nil {
		fmt.Println("err2:", err)
	}
	fmt.Println(content, valueType, offset)
	//jsonparser提供了解析bool、string、float64以及int64类型的方法，至于其他类型，我们可以通过valueType类型来自己进行转化
	result1, err := jsonparser.ParseString(content)
	if err != nil {
		fmt.Println("err3:", err)
	}
	fmt.Println(result1)

	err = jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		fmt.Printf("key:%s\n value:%s\n Type:%s\n", string(key), string(value), dataType)
		return nil
	}, "person", "avatars")

}
