package util

import (
	"github.com/patrickmn/go-cache"
	"yin/ip/util"
)
var Che *cache.Cache
const (
	Time_TIMEMSSQL string ="2006-01-02T15:04:05.999Z"
	Time_TIMEMYSQL string="2006-01-02T15:04:05+08:00"
)
//异常处理
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
func LogPrint(msg string,err error)  {
	if err!=nil{
		util.LogInfo().Println(msg,err)
	}
}