package main

import (
	"context"
	"fmt"
	"grpc-lb/api/protobuf-spec/template"
	"grpc-lb/internal/templateStore/service"
	"time"
)

func main() {

	for t := range time.NewTicker(1 * time.Second).C {
		client := service.NewTemplateClient()
		ctxc, _ := context.WithTimeout(context.Background(), 2*time.Second)
		resp, err := client.Get(ctxc, &template.ShowRequest{
			TemplateId: "T_08504EOI" + time.Now().Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			panic(err)
		}

		go func() {
			fmt.Println(resp.TemplateId + " : " + time.Since(t).String())
		}()

	}
}
