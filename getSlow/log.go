// @Time    :  2019/11/8
// @Software:  GoLand
// @File    :  log.go
// @Author  :  Abb1513

package getSlow

import (
	filename "github.com/keepeye/logrus-filename"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var logs = logrus.New()

func configLocalFilesystemLogger() {
	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	//baseLogPath := path.Join(".", "rdsSlowSql.log")
	writer, err := rotatelogs.New(
		//"rdsSlowSql" + "-%Y%m%d%H%M.log",
		"rdsSlowSql" + ".log",
	//	rotatelogs.WithLinkName(baseLogPath),        // 生成软链，指向最新日志文件
	//	rotatelogs.WithMaxAge(7*24*time.Hour),       // 文件最大保存时间
	//	rotatelogs.WithRotationTime(7*24*time.Hour), // 日志切割时间间隔
	)
	//logPath := lfshook.PathMap{logrus.}
	if err != nil {
		logs.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	// 为不同级别设置不同的输出目的
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{DisableTimestamp: false})
	logs.AddHook(filenameHook)
	logs.AddHook(lfHook)
}
