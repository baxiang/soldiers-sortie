package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

func main() {

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"}, // 集群列表
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 模拟etcd中kv的变化
	go func() {
		for {
			client.Put(context.TODO(), "hello", "world")

			client.Delete(context.TODO(), "hello")
			time.Sleep(2 * time.Second)
		}
	}()

	// 启动监听 5秒后关闭
	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(10*time.Second, func() {
		cancelFunc()
	})
	watchRespChan := client.Watch(ctx, "hello")

	// 处理kv变化事件
	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("put action", string(event.Kv.Value))
			case mvccpb.DELETE:
				fmt.Println("delete action", string(event.Kv.Key))
			}
		}
	}
}
