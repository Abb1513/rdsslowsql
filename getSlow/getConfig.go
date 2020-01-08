// @Time    :  2019/11/8
// @Software:  GoLand
// @File    :  getConfig.go
// @Author  :  Abb1513

package getSlow

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"time"
)

var V = viper.New()
var Conf Config
var db *gorm.DB
var err error

func init() {
	V.SetConfigName("config")
	V.AddConfigPath(".")
	V.SetConfigType("yml")
	if err := V.ReadInConfig(); err != nil {
		panic(err)
	}
	configLocalFilesystemLogger()
	GetConfig()
	db, err = gorm.Open("mysql", Conf.MysqlClient)
	if err != nil {
		logs.Panic("数据库连接失败, ", err)
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(10)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(100)

	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	db.DB().SetConnMaxLifetime(time.Hour)

}

type Config struct {
	DbInstanceId []string
	MysqlClient  string
}

func GetConfig() Config {
	// 获取配置文件

	if err := V.Unmarshal(&Conf); err != nil {
		logs.Error("读取配置文件, ", err)
	}
	logs.Info("获取到配置, ", Conf)
	return Conf
}
