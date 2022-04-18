package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
	   fmt.Printf("connect to etcd failed err,%v",err)
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	//str := `[{"path":"D:\\studyGo\\s4.log","topic":"s4_log"},{"path":"D:\\studyGo\\s5.log","topic":"s5_log"}]`
	//str := `[{"path":"D:\\studyGo\\s4.log","topic":"s4_log"},{"path":"D:\\studyGo\\s5.log","topic":"s5_log"},{"path":"D:\\studyGo\\s6.log","topic":"s6_log"}]`
	str := `[{"path":"D:\\studyGo\\s4.log","topic":"web_log"},{"path":"D:\\studyGo\\s5.log","topic":"s5_log"},{"path":"D:\\studyGo\\s6.log","topic":"s6_log"},{"path":"D:\\studyGo\\s7.log","topic":"s7_log"}]`
	_, err = cli.Put(ctx, "collect_log_192.168.31.217_conf", str)

	if err != nil {
		fmt.Printf("put to etcd failed,err:%v",err)
	}
	cancel()

    //get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	gr, err := cli.Get(ctx, "collect_log_192.168.31.217_conf")
	if err != nil {
		fmt.Printf("get to etcd failed,err:%v",err)
	}
	for _,ev := range gr.Kvs{
		fmt.Printf("key：%s value:%s\n",ev.Key,ev.Value)
	}
	cancel()
}