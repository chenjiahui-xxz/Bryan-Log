package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)
var (
	client sarama.SyncProducer
	// msgChan 之前是想MsgChan字符串，还是放成一个指针更省内存;
	msgChan chan *sarama.ProducerMessage
)

// Init 初始化kafka连接
func Init(address []string,chanSize int64) (err error){
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client,err = sarama.NewSyncProducer(address,config)
	if err != nil {
		logrus.Error("producer close,err:",err)
		return
	}
	//要使用全局的，需要先初始化,管道多少个写在配置文件里面
	msgChan = make(chan *sarama.ProducerMessage,chanSize)

	//起一个后台的goroutine从msgChan中读取数据
	go sendMsg()
    return
}

//从管道MsgChan中读取数据发到kafka
func sendMsg()  {
	for  {
		select {
		case msg := <-msgChan:
			pid,offset,err:=client.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msg err:",err)
				return
			}
			logrus.Infof("send msg to kafka success,pid:%v ,offset:%v",pid,offset)
		}
	}
}

// ToMsgChan  定义一个函数向外暴露msgchan
func ToMsgChan(msg  *sarama.ProducerMessage)  {
   msgChan<-msg
}