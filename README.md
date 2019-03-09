
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
    go run cmd/demo/svr/svr.go - port 50001
    go run cmd/demo/svr/svr.go - port 50002
    go run cmd/demo/svr/svr.go - port 50003

    # 启动客户端
    go run cmd/demo/cli/cli.go
    
## 使用配置文件
    cp config-eaxmple　config
    cd config
    mv config-eaxmple.go config.go
## 使用 redis mysql
    在　tool　下，有相应类
## 　使用 proto 编译出　go 服务端、客户端代码
    protoc --go_out=plugins=grpc:. hello.proto
    
## 获取可使用 grpc-server (推荐使用此类调用)
    var base = baseServer.NewServer("template_store_service")
    lis := base.GetAliveServer()
    
 ### todo list
    -1 Graceful shutdown 
    0 服务部署上线
    1 服务熔断  (可参考 go-kit + Hystrix)
    2 服务监控  (可参考 prometheus + alertmanager + grafana)
    3 服务降级
    4 node + php 完整 client　demo