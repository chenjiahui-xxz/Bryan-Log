package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"log_transfer/es"
)

//初始化kafka。同时取出数据

func Init(addr []string,topic string)(err error){

	// 生成消费者 实例
	consumer, err := sarama.NewConsumer(addr, nil)
	if err != nil {
		log.Printf("init newConsumer fail,err:%v",err)
		return
	}
	// 拿到 对应主题下所有分区
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Printf("init Partitions fail,err:%v",err)
		return
	}
	//fmt.Println(topic)
    //fmt.Println(addr)
	for partition := range partitionList{
		//消费者 消费 对应主题的 具体 分区 指定 主题 分区 offset  return 对应分区的对象
		var pc sarama.PartitionConsumer
		pc, err = consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("consumer.ConsumePartition fail ,err:%v\n",err)
			return
		}

		// 运行完毕记得关闭
		//defer pc.AsyncClose()

		// 通过异步 拿到 消息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages(){
				//logDataChan <- msg //将取出来的日志数据给出管道
				var m1 map[string]interface{}
				err := json.Unmarshal(msg.Value, &m1)
				if err != nil {
					log.Printf("Unmarshal log msg fail ,err:%v\n",err)
					continue
				}
				es.PutLogData(m1)
			}
		}(pc)
	}
	return
}