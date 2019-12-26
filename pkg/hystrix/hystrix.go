package hystrix

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	_ "github.com/afex/hystrix-go/hystrix"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

func init() {
	hystrix.ConfigureCommand("seckill", hystrix.CommandConfig{
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
