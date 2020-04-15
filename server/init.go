package server

import (
	"errors"
	"fmt"
	"log"
	"reptile-go/model"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var DbEngin *xorm.Engine

func init() {
	// 读取yaml配置文件
	viper.SetConfigName("config1")
	// 设置配置文件的搜索目录
	viper.AddConfigPath("./config/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			log.Println("read config error")
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}
	username := viper.GetString(`mysql.username`)
	password := viper.GetString(`mysql.password`)
	host := viper.GetString(`mysql.host`)
	port := viper.GetString(`mysql.port`)
	dbname := viper.GetString(`mysql.dbname`)
	DsName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)
	driveName := "mysql"
	err := errors.New("")
	DbEngin, err = xorm.NewEngine(driveName, DsName)
	if nil != err && "" != err.Error() {
		panic(err.Error())
	}
	// 是否显示SQL语句
	DbEngin.ShowSQL(true)
	// 设置数据库最大打开的连接数
	DbEngin.SetMaxOpenConns(2)
	// 自动建表
	DbEngin.Sync2(new(model.User),
		new(model.Contact),
		new(model.Community),
		new(model.Message))
	fmt.Println("init data base ok")
}
