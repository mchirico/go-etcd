package main

import (  
    "fmt"
    "context"
    "log"
    "github.com/etcd-io/etcd/clientv3"
    "time"
    "strconv"
)

var (  
    dialTimeout    = 2 * time.Second
    requestTimeout = 10 * time.Second
)

func main() {  
    ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
    cli, _ := clientv3.New(clientv3.Config{
        DialTimeout: dialTimeout,
        Endpoints: []string{"127.0.0.1:2379"},
    })
    defer cli.Close()
    kv := clientv3.NewKV(cli)

    _ = kv
}
