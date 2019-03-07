
# 测试
    govendor add +external  添加依赖包至 VENDOR下
## 启动ETCD
``` 
ubuntu:
apt install etcd
centos:
yum install etcd
mac:
brew install etcd


etcd  将会启动在：localhost:2379
```
## 启动测试程序

    # 分别启动服务端
    go run cmd/svr/svr.go - port 50001
    go run cmd/svr/svr.go - port 50002
    go run cmd/svr/svr.go - port 50003

    # 启动客户端
    go run cmd/cli/cli.go
    
## 使用 redis mysql
    在　tool　下，有相应类
## 使用 baseGrpc
    
	var base = basegrpc.InitGrpc{ServiceName: "site"}
	cli := base.NewBaseGrpc()

	s := grpc.NewServer()
	site.RegisterSiteServer(s, &Server{})

	s.Serve(cli)
## 　使用 proto 编译出　go 服务端、客户端代码
    protoc --go_out=plugins=grpc:. hello.proto