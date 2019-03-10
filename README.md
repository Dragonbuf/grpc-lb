# 测试
    govendor add +external  添加依赖包至 VENDOR下
## 启动ETCD
``` 
ubuntu:
    1、安装
ETCD_VER=v3.1.0
DOWNLOAD_URL=https://github.com/coreos/etcd/releases/download
curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
mkdir -p /tmp/test-etcd && tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/test-etcd --strip-components=1
/tmp/test-etcd/etcd --version
/tmp/test-etcd/etcd
2、增加查看keys
sudo apt install httpie
http PUT http://127.0.0.1:2379/v2/keys/message value=="hello, etcd"
或者
curl -L http://127.0.0.1:2379/v2/keys/mykey -XPUT -d value="this is awesome"
登录http://localhost：2379/keys查看keys
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