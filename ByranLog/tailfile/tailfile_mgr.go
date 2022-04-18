package tailfile

import (
	"ByranLog/common"
	"github.com/sirupsen/logrus"
)

//管理tailTask，不是全局，但是是透传下去的，也还是同一个
type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask        //所有的tailTask任务，也就是任务列表
	collectEntryList []common.CollectEntity      //所有的配置项。也就是读到的新的配置列表，main函数第一次启动的时候就是放在这边
	confChan         chan []common.CollectEntity //等待新的配置出现
}

var (
	ttMgr *tailTaskMgr
)

// Init 主要是为了在main函数调用启动
func Init(allConf []common.CollectEntity) (err error) {
	//多个任务在这边，循环每一个任务，在这边开个协程开始去读取数据--这是第一次启动得时候读取得配置
	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20),
		collectEntryList: allConf,
		confChan:         make(chan []common.CollectEntity),
	}
	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic)
		err = tt.Init()
		if err != nil {
			logrus.Error("tailfile: create tailObj for path:%s failed,err:%v\n", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success", conf.Path)

		ttMgr.tailTaskMap[tt.path] = tt //已经有的任务保存起来

		go tt.run() //去收集日志
	}

	//需要开个协程去监听配置有没有改变，因为etcdWatch会监听到数据改变发到管道中。这边去读取管道中得数据即可
	go ttMgr.watch() //ttMgr这里开始透传

	return
}

//自己开个协程在读取外面etcd送进来的数据
func (t *tailTaskMgr) watch() {
	for {
		newConf := <-t.confChan //就是etcd监控送进来的数据
		logrus.Infof("get new conf from etcd,conf:%v，start manage TailTask...", newConf)

		for _, conf := range newConf {
			//1.原来已经存在的任务不用动
			if t.isExist(conf) {
				continue
			}
			//2.原来没有的我要新起一个任务
			tt := newTailTask(conf.Path, conf.Topic)
			err := tt.Init() //这两句就是新起任务

			if err != nil {
				logrus.Error("tailfile: create tailObj for path:%s failed,err:%v\n", conf.Path, err)
				continue
			}
			logrus.Infof("create a tail task for path:%s success", conf.Path)

			t.tailTaskMap[tt.path] = tt //已经有的任务保存起来
			go tt.run()                 //任务开始执行哈哈
		}
		//3.原来有，现在没有需要停掉
		//找出任务在跑的存在的，但是现在来的配置是没有的，将他们停掉
		for key, task := range t.tailTaskMap {
			var found bool
			for _, conf := range newConf {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				//找不到证明就是已经不存在了，要把他停掉
				logrus.Infof("the task collect path:%s need to stop", task.path)
				delete(t.tailTaskMap, key) //从管理的任务列表中删除
				task.cancel()
			}
		}

	}
}

func (t *tailTaskMgr) isExist(conf common.CollectEntity) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok
}

// SendNewConf 给外面暴露塞进去管道的
func SendNewConf(newConf []common.CollectEntity) {
	ttMgr.confChan <- newConf
}
