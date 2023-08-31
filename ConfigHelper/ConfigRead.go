package config

import (
	"io/ioutil"

	"github.com/antonholmquist/jason"
)

var dir string = "./ConfigHelper/config.json"

func ReadValue(valKey string) string {
	var data []byte
	data, _ = ioutil.ReadFile(dir)
	v, _ := jason.NewObjectFromBytes(data)
	rtrnVal, _ := v.GetString(valKey)
	return rtrnVal
}
