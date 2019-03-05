package main

import (
	"fmt"
	_ "github.com/coreos/etcd/mvcc/mvccpb"
	"grpc-lb/cmd/etcd-web-watch"
	_ "net/http"
	_ "strings"
)

func main() {
	etcd_web_watch.Start()
}

func UploadConfig() {
	config := "/config/grpc/db"
	fmt.Println(config)
}
