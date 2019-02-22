package main

import (
	"context"
	"google.golang.org/grpc"
	"grpc-lb/cmd/basegrpc"
	"grpc-lb/cmd/server/site/proto"
	"grpc-lb/tool"
)

func main() {
	var base = basegrpc.InitGrpc{ServiceName: "site"}
	cli := base.NewBaseGrpc()

	s := grpc.NewServer()
	site.RegisterSiteServer(s, &Server{})

	_ = s.Serve(cli)
}

type Server struct{}

func (s *Server) Show(ctx context.Context, in *site.ShowRequest) (out *site.ShowReply, err error) {

	out.Title = in.GetUrl()
	out.Reset()

	redis := tool.RedisPool.Get()
	redis.Do("set", "key", "value")

	return out, nil
}
