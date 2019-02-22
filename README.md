
# 测试
    govendor add +external
## 启动ETCD
ubuntu:
apt install etcd
centos:
yum install etcd
mac:
brew install etcd


etcd localhost:2379
## 启动测试程序

    # 分别启动服务端
    go run cmd/svr/svr.go - port 50001
    go run cmd/svr/svr.go - port 50002
    go run cmd/svr/svr.go - port 50003

    # 启动客户端
    go run cmd/cli/cli.go
