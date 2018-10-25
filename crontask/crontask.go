package crontask

import (
	"time"

	"github.com/hexiaoyun128/gin-base-framework/utils"
)

const (
	// check and create log 30 second
	checkLog = 30 * time.Second
)

//定时任务
func CronTask() {
	go func() {
		logTicker := time.NewTicker(checkLog)

		defer func() {
			logTicker.Stop()
		}()

		for utils.Run {
			select {
			case <-logTicker.C:
				// 定时处理log日志
				utils.CheckLog()
			}
		}
	}()
}
