package main

import (
	"ByranLog/common"
	"ByranLog/etcd"
	"ByranLog/kafka"
	"ByranLog/tailfile"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type Config struct {
	KafkaConfig KafkaConfig `ini:"kafka"`
	Collect     Collect     `ini:"collect"`
	EtcdConfig  EtcdConfig  `ini:"etcd"`
}
type KafkaConfig struct {
	Address  string `ini:"address"`
	ChanSize int64  `ini:"chan_size"`
}
type Collect struct {
	LogFilePath string `ini:"logfile_path"`
}
type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

func main() {
	//-1：获取本机ip
	ip, err := common.GetIp()
	if err != nil {
		logrus.Errorf("get ip failed,err: %v", err)
		return
	}

	var configObj = new(Config) //new一个指针
	//0.读配置文件
	err = ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Errorf("Fail to read file: %v", err)
		return
	}
	fmt.Printf("%#v\n", configObj)
	//1.初始化连接kafka
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka err:%v", err)
		return
	}
	logrus.Info("init kafka success!")

	//初始化etcd连接
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd err:%v", err)
		return
	}
	//从etcd中拉取要收集日志的配置项,现在根据本机ip
	collectKey := fmt.Sprintf(configObj.EtcdConfig.CollectKey, ip)
	aliConf, err := etcd.GetConf(collectKey)
	if err != nil {
		logrus.Errorf("Get Conf from etcd failed,err:%v", err)
		return
	}
	fmt.Println(aliConf)
	//去监控etcd中的key即是configObj.EtcdConfig.CollectKey对应值的变化
	go etcd.WatchConf(collectKey)
	//2.根据配置中的日志路径使用tail去收集日志
	err = tailfile.Init(aliConf) //将读取到的etcd配置去生成读取文件任务
	if err != nil {
		logrus.Errorf("init tailfile err:%v", err)
		return
	}
	logrus.Info("init tailfile success!")
	run()
}

func run() {
	for {
		select {}
	}
}
