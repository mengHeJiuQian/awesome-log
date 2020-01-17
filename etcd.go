package main

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/clientv3"
	"time"
)

/**
 * Create Time : 2020/1/15 下午3:05
 * Update Time :
 * Author : sheldon
 * Decription :
 */

var Client *clientv3.Client
var logConfChan chan string

// 初始化etcd
func initEtcd(addr []string, keyfmt string, timeout time.Duration) (err error) {

	var keys []string
	for _, ip := range ipArrays {
		//keyfmt = /logagent/%s/log_config
		fmt.Println("aaa")
		fmt.Println(keyfmt)
		fmt.Println(ip)

		keys = append(keys, fmt.Sprintf(keyfmt, ip))
	}
	logs.Info("keys=%v", keys)

	logConfChan = make(chan string, 10)
	logs.Debug("etcd watch key:%v timeout:%v", keys, timeout)

	Client, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: timeout,
	})
	defer Client.Close() // ++
	if err != nil {
		logs.Error("connect failed,err:%v", err)
		return
	}
	logs.Debug("init etcd success")
	waitGroup.Add(1)
	for _, key := range keys {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		//++
		_, err = Client.Put(ctx, key, key)
		if err != nil {
			logs.Error("etcd存值错误，key=%s, err=%v", key, err)
		}
		//++

		// 从etcd中获取要收集日志的信息
		resp, err := Client.Get(ctx, key)
		cancel()
		if err != nil {
			logs.Warn("get key %s failed,err:%v", key, err)
			continue
		}

		for _, ev := range resp.Kvs {
			logs.Debug("%q : %q\n", ev.Key, ev.Value)
			logConfChan <- string(ev.Value)
		}
	}
	go WatchEtcd(keys)
	return
}

func WatchEtcd(keys []string) {
	// 这里用于检测当需要收集的日志信息更改时及时更新
	var watchChans []clientv3.WatchChan

	logs.Info("keys=%v", keys)

	for _, key := range keys {
		rch := Client.Watch(context.Background(), key)
		watchChans = append(watchChans, rch)
	}

	for {
		for _, watchC := range watchChans {
			select {
			case wresp := <-watchC:
				for _, ev := range wresp.Events {
					logs.Debug("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					logConfChan <- string(ev.Kv.Value)
				}
			default:

			}
		}
		time.Sleep(time.Second)
		break
	}
	waitGroup.Done()
}

func GetLogConf() chan string {
	return logConfChan
}
