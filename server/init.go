package server

import (
	"errors"
	"fmt"
	"reptile-go/config"
	"reptile-go/model"

	"github.com/importcjj/sensitive"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var DbEngin *xorm.Engine
var vip *viper.Viper

var ViperConfig config.Configuration

// 初始化配置信息
func init() {
	runtimeViper := viper.New()
	runtimeViper.AddConfigPath("./config/") // 设置配置文件的搜索目录
	runtimeViper.SetConfigName("config")    // 配置文件名
	runtimeViper.SetConfigType("yaml")      // 文件格式
	if err := runtimeViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	runtimeViper.Unmarshal(&ViperConfig)
	//	监听配置文件变更
	runtimeViper.WatchConfig()
	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		runtimeViper.Unmarshal(&ViperConfig)
		//fmt.Println("变化：" + ViperConfig.MySQL.Database)
		//ViperConfig.LocaleBundle.MustLoadMessageFile(ViperConfig.App.Locale+"/active."+ViperConfig.App.Language+".json")
	})
}

// 初始化数据库
func init() {
	username := ViperConfig.MySQL.Username
	password := ViperConfig.MySQL.Password
	host := ViperConfig.MySQL.Address
	port := ViperConfig.MySQL.Port
	dbname := ViperConfig.MySQL.Database
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

// 敏感词典
var Filter *sensitive.Filter

func init() {
	Filter = sensitive.New()
	err := Filter.LoadWordDict("./dict/sensitive.txt")
	if err != nil {
		fmt.Println(err)
	}
}
