### 介绍：使用 go 1.13 以上版本
个人研究的 go grpc 微服务框架

#### 启动项目　（3 步即可）
> 设置环境变量，防止被墙：export GOPROXY=http://mirrors.aliyun.com/goproxy（win: $env:GOPROXY = "http://mirrors.aliyun.com/goproxy" 或 https://goproxy.io）

> cp -R config-example config (需要修改 config 下数据库、redis 密码)

> go run ./cmd/template/server.go

### ubuntu install protoc
    1 下载 https://github.com/protocolbuffers/protobuf/releases 相应版本

    2 解压后，复制 protoc 至 /usr/local/bin/下

    3 尝试 protoc --version 是否成功

    4 go get  github.com/golang/protobuf/protoc-gen-go

    5 go install github.com/golang/protobuf/protoc-gen-go

    6 cd /protos/demo 运行下列命令既可编译出　go 服务端、客户端代码
        protoc --go_out=plugins=grpc:. hello.proto

```
### todo list
    -1 Graceful shutdown 
    0 服务部署上线
    1 服务熔断  (可参考 go-kit + Hystrix)
    2 服务监控  (可参考 prometheus + alertmanager + grafana)
    3 服务降级