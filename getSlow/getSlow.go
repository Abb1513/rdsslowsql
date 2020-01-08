// @Time    :  2019/11/8
// @Software:  GoLand
// @File    :  getSlow.go
// @Author  :  Abb1513

package getSlow

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
)

type sqlSlowRecord struct {
	rds.SQLSlowRecord
	DbInstanceId string
}

func AliApi(startTime string, endTime string, dbInstanceId string, pageNum int) (int, int) {
	// 请求aliyun API 获取慢日志
	logs.Infof("收到id: %s, 开始时间startTime: %s, 结束时间endTime: %s", dbInstanceId, startTime, endTime)
	accessKeyId := ""
	accessSecret := ""
	client, err := rds.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessSecret)
	if err != nil {
		logs.Panic("连接阿里云失败, ", err)
	}
	request := rds.CreateDescribeSlowLogRecordsRequest()
	request.Scheme = "https"
	request.DBInstanceId = dbInstanceId
	request.StartTime = startTime
	request.EndTime = endTime
	request.PageSize = requests.NewInteger(100)
	request.PageNumber = requests.NewInteger(pageNum)
	response, err := client.DescribeSlowLogRecords(request)
	if err != nil {
		logs.Panic("DescribeSlowLogRecords, 返回, ", err)
	}

	logs.Infof("总数:%d ,第%d页,本页%d条 ,收到返回:%s ", response.TotalRecordCount, response.PageNumber, response.PageRecordCount, response.Items.SQLSlowRecord)
	sqlSlowRecordList := response.Items.SQLSlowRecord

	//logs.Info("收到返回,sqlSlowList, ", sqlSlowList)
	var sql sqlSlowRecord
	if len(sqlSlowRecordList) > 0 {
		for _, k := range sqlSlowRecordList {
			sql.DbInstanceId = dbInstanceId
			sql.DBName = k.DBName
			sql.ExecutionStartTime = toGMT(k.ExecutionStartTime)
			sql.HostAddress = k.HostAddress
			sql.LockTimes = k.LockTimes
			sql.ParseRowCounts = k.ParseRowCounts
			sql.SQLText = k.SQLText
			sql.ReturnRowCounts = k.ReturnRowCounts
			sql.QueryTimes = k.QueryTimes
			writeMysql(sql, dbInstanceId)
		}
	} else {
		logs.Info("没有慢日志")
	}
	return response.TotalRecordCount, response.PageRecordCount
}

func GetUTCTime() (startTime, endTime string) {
	// 生成固定格式时间
	now := time.Now()
	sd, _ := time.ParseDuration("-24h")
	s, _ := time.ParseDuration("-8h")
	endTime = now.Add(s).Format("2006-01-02T15:04Z")
	startTime = now.Add(sd).Format("2006-01-02T15:04Z")
	//logs.Info(time)
	return
}

func toGMT(times string) string {
	// 转换时间到 北京时间
	timeLayoutStr := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	gmtTime, _ := time.ParseInLocation("2006-01-02T15:04:05Z", times, loc)
	exTime := gmtTime.Format(timeLayoutStr)
	return exTime
}

func writeMysql(s sqlSlowRecord, id string) {
	// 持久化存储
	db.Table("sql_slow_record").Create(&s)
	//db.Create(&s)
	if res := db.Table("sql_slow_record").NewRecord(&s); res {
		logs.Info("写入成功, ", s)
	}
}
