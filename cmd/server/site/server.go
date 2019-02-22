package main

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
	"grpc-lb/cmd/basegrpc"
	"grpc-lb/cmd/server/site/proto"
	"grpc-lb/tool"
)

func main() {

	client := tool.RedisPool.Get()

	if client.Err() != nil {
		fmt.Println(client.Err())
	}

	res, err := client.Do("set", "key", "value")

	if err == nil {
		fmt.Println(res)
	}

	redis2 := tool.RedisPool.Get()
	res, err = redis.String(redis2.Do("get", "key"))
	if err == nil {
		fmt.Println(res)
	}

	var base = basegrpc.InitGrpc{ServiceName: "site"}
	cli := base.NewBaseGrpc()

	s := grpc.NewServer()
	site.RegisterSiteServer(s, &Server{})

	err = s.Serve(cli)
	if err != nil {
		panic(err)
	}

}

type Server struct{}

func (s *Server) Show(ctx context.Context, in *site.ShowRequest) (out *site.ShowReply, err error) {

	out.Title = in.GetUrl()
	out.Reset()

	redis := tool.RedisPool.Get()

	if redis.Err() != nil {
		fmt.Println(redis.Err())
	}

	res, err := redis.Do("set", "key", "value")

	if err != nil {
		fmt.Println(res)
	}

	redis2 := tool.RedisPool.Get()
	redis2.Do("get", "key")

	//redis3 := tool.RedisPool.Get()
	//redis3.Do("del", "key")

	return out, nil
}
