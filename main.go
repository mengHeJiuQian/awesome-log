package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

/**
 * Create Time : 2020/1/15 下午3:07
 * Update Time :
 * Author : sheldon
 * Decription :
 */

func getLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "trace":
		return logs.LevelTrace
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "error":
		return logs.LevelError
	default:
		return logs.LevelDebug
	}
}

func initLog() (err error) {
	//初始化日志库
	config := make(map[string]interface{})
	config["filename"] = "./logs/logcollect.log"
	config["level"] = getLevel(appConfig.LogLevel) // main函数调用appConfig()就初始化过了appConfig全局变量
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("mashal failed,err:", err)
		return
	}
	logs.SetLogger(logs.AdapterConsole, string(configStr))
	logs.Info("日志配置初始化成功。initLog()")
	return
}

func main() {
	err := initConfig("./conf/app.conf")
	if err != nil {
		panic(fmt.Sprintf("init config failed,err:%v\n", err))
	}
	err = initLog()
	if err != nil {
		logs.Info("初始化日志出错")
		return
	}
	logs.Debug("init success")
	ipArrays, err = getLocalIP()
	logs.Info(ipArrays)
	if err != nil {
		logs.Error("get local ip failed, err:%v", err)
		return
	}

	logs.Debug("get local ip succ, ips:%v", ipArrays)
	err = initKafka()
	if err != nil {
		logs.Error("init kafka faild, err:%v", err)
		return
	}
	err = initEtcd(appConfig.etcdAddr, appConfig.etcdWatchKeyFmt, time.Duration(appConfig.etcdTimeout)*time.Millisecond)
	if err != nil {
		logs.Error("init etcd failed, err:%v", err)
		return
	}
	RunServer()
}
