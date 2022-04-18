package tailfile

import (
	"ByranLog/kafka"
	"context"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

//tail 相关
type tailTask struct {
	path   string
	topic  string
	tObj   *tail.Tail
	ctx    context.Context
	cancel context.CancelFunc
}

//构造函数
func newTailTask(path, topic string) *tailTask {
	ctx, cancel := context.WithCancel(context.Background())
	tt := &tailTask{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}
	return tt
}

// 实行完构造函数后就初始化tObj这个读取数据的句柄
func (t *tailTask) Init() (err error) {
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.tObj, err = tail.TailFile(t.path, cfg)
	return
}

func (t *tailTask) run() {
	//读取日志发送到管道
	logrus.Infof("collect path:%s is running..", t.path)
	for {
		select {
		case <-t.ctx.Done(): //如果外面ctx调用了cancel函数，这边就直接结束了
			logrus.Infof("path:%s is stopping...", t.path)
			return
		case line, ok := <-t.tObj.Lines:
			//循环去监控这个日志，读取出来
			if !ok {
				logrus.Error("tails file close reopen,filename: %s\n", t.tObj.Filename)
				time.Sleep(time.Second)
				continue
			}
			//如果是空行，不要采集
			//fmt.Printf("#%v",line.Text)
			if len(strings.Trim(line.Text, "\r")) == 0 {
				logrus.Info("出现空行，直接跳过..")
				continue
			}
			//fmt.Println(line.Text)
			//如果直接发的话，一个一个就往kafka，那么如果慢的话这边就会阻塞掉，还是利用管道异步去执行
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic
			msg.Value = sarama.StringEncoder(line.Text)
			//丢到管道中
			kafka.ToMsgChan(msg)
		}
	}
}
