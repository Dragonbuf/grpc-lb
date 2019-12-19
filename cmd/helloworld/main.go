package main

import (
	"context"
	"fmt"
	templateStore "grpc-lb/internal/app/server/templateStore/proto"
	"grpc-lb/internal/pkg/loadBalance"
	"time"
)

func main() {

	roundrobinConn, err := loadBalance.NewBaseClient("template_store_service").GetRoundRobinConn()
	if err != nil {
		panic(err)
	}

	defer roundrobinConn.Close()

	ticker := time.NewTicker(1000 * time.Millisecond)
	for t := range ticker.C {
		ctxc, _ := context.WithTimeout(context.Background(), 2*time.Second)
		client := templateStore.NewTemplateStoreClient(roundrobinConn)

		resp, _ := client.Get(ctxc, &templateStore.ShowRequest{
			TemplateId:           "T_49TZW9W7",
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		})

		go func() {
			fmt.Print("i got template:" + resp.TemplateId + "in ")
			fmt.Println(time.Since(t))
		}()

	}
}
