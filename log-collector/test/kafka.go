package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

/**
 * Create Time : 2020/1/11 上午11:33
 * Update Time :
 * Author : sheldon
 * Decription : 对kafka的操作，包括初始化kafka连接，以及给kafka发送消息
 */

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	// 消息信息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "javaWeb_log"
	msg.Value = sarama.StringEncoder("this is a good test,my message is zhaofan")

	client, err := sarama.NewSyncProducer([]string{"aliyun:9092"}, config)

	if err != nil{
		fmt.Println("producer close err:",err)
		return
	}

	// 1. 可以优化到创建之后就执行
	defer client.Close()

	// 没两秒生产一条消息
	for {
		pid, offset, err := client.SendMessage(msg)
		if err != nil {
			fmt.Println("send message failed,",err)
			return
		}
		fmt.Printf("pid:%v offset:%v\n",pid,offset)
		time.Sleep(2*time.Second)
	}

	/* // 每次生产一条消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send message failed,",err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n",pid,offset)
	 */
}
