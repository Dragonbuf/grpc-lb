package main

import (
	hystrix2 "grpc-lb/pkg/hystrix"
	"log"
	"testing"
)

func TestDO(t *testing.T) {
	for i := 0; i < 10; i++ {
		err := hystrix2.DoTest()
		if err != nil {
			log.Printf("testHystrix error:%v", err)
		}
	}
}
