package etcd

import (
	"ByranLog/common"
	"ByranLog/tailfile"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client/v3"
	"time"
)

var (
	client *clientv3.Client
)

func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed err,%v", err)
	}
	return
}

// GetConf 拉取日志收集的配置项目
func GetConf(key string) (collectEntityList []common.CollectEntity, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get etcd failed,err:%v", err)
		return
	}

	if len(resp.Kvs) == 0 {
		logrus.Warningf("get len:0 conf from etcd by key :%s", key)
		return
	}

	ret := resp.Kvs[0]

	err = json.Unmarshal(ret.Value, &collectEntityList)

	if err != nil {
		logrus.Errorf("json Unmarshal failed,err:%v", err)
		return
	}
	return
}

// WatchConf 监控etcd中日志收集项配置的变化
func WatchConf(key string) {
	for {
		watchCh := client.Watch(context.Background(), key)

		for wries := range watchCh {
			logrus.Infof("get new conf from etcd!")
			for _, ev := range wries.Events {
				fmt.Printf("type:%s %q:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				var newConf []common.CollectEntity
				if ev.Type == clientv3.EventTypeDelete {
					//如果是删除
					logrus.Warning("坑爹玩意,etcd delete the key !!")
					tailfile.SendNewConf(newConf)
					continue
				}
				err := json.Unmarshal(ev.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("json Unmarshal new conf failed,err:%v", err)
					continue
				}
				//应该使用新的配置去生成读取文件日志，那么需要告知taiFile中得init中的方法去处理
				tailfile.SendNewConf(newConf) //因为这个管道是无缓存的，没有人在这边接收就是阻塞
			}
		}
	}
}
