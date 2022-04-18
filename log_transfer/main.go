package main

import (
	"fmt"
	"log_transfer/es"
	"log_transfer/kafka"
	"log_transfer/model"
)
import "gopkg.in/ini.v1"

//写一个消费者，就是将kafka的配置文件读到es里面，在从kiban中显示出来

func main() {
	//1.加载配置文件
	var cfg = new(model.Config)
	err := ini.MapTo(cfg, "./config/logtransfer.ini")
	if err != nil {
		fmt.Printf("load config fail,err:%v\n", err)
		panic(err)
	}
	fmt.Println("load config success ")

	//2.连接es
	err = es.Init(cfg.ESConfig.Address, cfg.ESConfig.Index, cfg.ESConfig.GoroutineNum, cfg.ESConfig.MaxSize)
	if err != nil {
		fmt.Printf("init es fail,err:%v\n", err)
		panic(err)
	}
	fmt.Println("connect es success ")

	//3.连接kafka
	err = kafka.Init([]string{cfg.KafkaConfig.Address}, cfg.KafkaConfig.Topic)
	if err != nil {
		fmt.Printf("connect kafka fail,err:%v\n", err)
		panic(err)
	}
	fmt.Println("connect kafka success ")

    //停顿
	select {}
}


