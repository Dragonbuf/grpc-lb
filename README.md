### 项目目录
```
/cmd
main函数文件（比如 /cmd/myapp.go）目录，这个目录下面，每个文件在编译之后都会生成一个可执行的文件。

不要把很多的代码放到这个目录下面，这里面的代码尽可能简单。

/internal
应用程序的封装的代码，某个应用私有的代码放到 /internal/myapp/ 目录下，多个应用通用的公共的代码，放到 /internal/common 之类的目录。

/pkg
一些通用的可以被其他项目所使用的代码，放到这个目录下面

/vendor
项目依赖的其他第三方库，使用 glide 工具来管理依赖

/api
协议文件，Swagger/thrift/protobuf 等

/web
web服务所需要的静态文件

/configs
配置文件

/init
服务启停脚本

/scripts
其他一些脚本，编译、安装、测试、分析等等

/build
持续集成目录

云 (AMI), 容器 (Docker), 操作系统 (deb, rpm, pkg)等的包配置和脚本放到 /build/package/ 目录

/deployments
部署相关的配置文件和模板

/test
其他测试目录，功能测试，性能测试等

/docs
设计文档

/tools
常用的工具和脚本，可以引用 /internal 或者 /pkg 里面的库

/examples
应用程序或者公共库使用的一些例子

/assets
其他一些依赖的静态资源
```

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