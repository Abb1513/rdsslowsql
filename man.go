// @Time    :  2019/11/8
// @Software:  GoLand
// @File    :  man.go
// @Author  :  Abb1513

package main

import (
	"sqlslow/getSlow"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	// 定时任务 每天7点30
	c.AddFunc("CRON_TZ=Asia/Shanghai 50 7 * * *", func() {
		startTime, endTime := getSlow.GetUTCTime() // 开始时间
		Conf := getSlow.GetConfig()
		for _, v := range Conf.DbInstanceId {
			pageNum := 1
			var totalNum int
			total, count := getSlow.AliApi(startTime, endTime, v, pageNum)
			totalNum += count
			for totalNum != total {
				pageNum++
				_, count := getSlow.AliApi(startTime, endTime, v, pageNum)
				totalNum += count
			}
		}

	})
	go c.Start()
	defer c.Stop()
	// 阻塞
	select {}
}
