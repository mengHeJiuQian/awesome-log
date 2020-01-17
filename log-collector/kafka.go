package main

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

/**
 * Create Time : 2020/1/11 上午11:33
 * Update Time :
 * Author : sheldon
 * Decription : 对kafka的操作，包括初始化kafka连接，以及给kafka发送消息
 */

var (
	client      sarama.SyncProducer
	kafkaSender *KafkaSender
)

type Message struct {
	line  string
	topic string
}

type KafkaSender struct {
	client   sarama.SyncProducer
	lineChan chan *Message
}

// 初始化kafka
func NewKafkaSender(kafkaAddr string) (kafka *KafkaSender, err error) {
	kafka = &KafkaSender{
		lineChan: make(chan *Message, 10000),
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		logs.Error("init kafka client failed,err:%v\n", err)
		return
	}
	kafka.client = client
	for i := 0; i < appConfig.KafkaThreadNum; i++ {
		// 根据配置文件循环开启线程去发消息到kafka
		go kafka.sendToKafka(i + 1)
	}
	return
}

func initKafka() (err error) {
	kafkaSender, err = NewKafkaSender(appConfig.kafkaAddr)
	return
}

func (k *KafkaSender) sendToKafka(kafkaNo int) {
	//从channel中读取日志内容放到kafka消息队列中
	for v := range k.lineChan {
		msg := &sarama.ProducerMessage{}
		msg.Topic = v.topic
		msg.Value = sarama.StringEncoder(v.line)
		_, _, err := k.client.SendMessage(msg)
		if err != nil {
			logs.Error("send message to kafka failed,err:%v", err)
		}
	}
	logs.Info("==>kafka生产者%v执行结束,协程退出", kafkaNo)
}

func (k *KafkaSender) addMessage(line string, topic string) (err error) {
	//我们通过tailf读取的日志文件内容先放到channel里面
	k.lineChan <- &Message{line: line, topic: topic}
	return
}

/*
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

	// 每次生产一条消息
	//pid, offset, err := client.SendMessage(msg)
	//if err != nil {
	//	fmt.Println("send message failed,",err)
	//	return
	//}
	//fmt.Printf("pid:%v offset:%v\n",pid,offset)

} */
