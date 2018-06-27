package main

import (
    "yin/xxd/crontask"
    "yin/xxd/hyperttp/server"
    "yin/xxd/util"
    "yin/xxd/ttd"
    "github.com/patrickmn/go-cache"
    "time"
)
func main() {
    util.Che=cache.New(10*time.Minute, 15*time.Minute)
    crontask.CronTask()
    go server.InitHttp()
    go ttd.InitWs()
    exitServer()
}
func exitServer() {
    for util.Run && util.GetNumGoroutine() > 2 {
        //util.Println("sleep ...")
        util.Sleep(3)
    }
}