package hystrix

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/pkg/errors"
	"grpc-lb/api/protobuf-spec/template"
	"grpc-lb/internal/template/service"
	"math/rand"
	"time"
)

var Name = "seckill"

func init() {
	hystrix.ConfigureCommand(Name, hystrix.CommandConfig{
		Timeout:                100, //cmd的超时时间，一旦超时则返回失败
		MaxConcurrentRequests:  1,   //最大并发请求数
		RequestVolumeThreshold: 5,   //熔断探测前的调用次数
		SleepWindow:            1,   //熔断发生后的等待恢复时间
		ErrorPercentThreshold:  10,  //失败占比
	})
}

func DoTest() error {
	query := func() error {
		fmt.Println("id got query")
		var err error
		r := rand.Float64()
		if r < 0 {
			err = errors.New("bad luck")
			return err
		} else {
			time.Sleep(20 * time.Millisecond)
		}

		return nil
	}

	service := "seckill"
	var err error
	err = hystrix.Do(service, func() error {
		err = query()
		return err
	}, nil)

	return err
}

func GetTemplateGetMiddle() service.Middle {
	return func(ctx context.Context, in interface{}) (interface{}, error) {
		if req, ok := in.(*template.ShowRequest); ok {
			req.TemplateId = "i am get template get middle : i got you,templateId"
		}
		return nil, nil
	}
}
