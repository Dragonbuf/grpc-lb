package etcd_web_watch

import (
	"fmt"
	_ "github.com/coreos/etcd/mvcc/mvccpb"
	_ "net/http"
	"strconv"
	_ "strings"
	"testing"
)

func TestStart(t *testing.T) {
	Start()
}

func TestServiceMap_AddMap(t *testing.T) {
	svrMap.AddMap("shit", "127.0.0.1:9000")
	fmt.Println(svrMap)

	if svrMap.svrMap["shit"].TotalService != 1 {
		t.Error("测试失败")
	}

	svrMap.AddMap("shit", "127.0.0.1:9001")
	if svrMap.svrMap["shit"].TotalService != 2 {
		t.Error("测试失败")
	}
}

func BenchmarkServiceMap_AddMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		svrMap.AddMap("i", strconv.Itoa(i))
	}
}

func BenchmarkIndexHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := i
		a++
	}
}

func TestServiceMap_DelMap(t *testing.T) {
	svrMap.AddMap("shit", "127.0.0.1:9000")
	svrMap.DelMap("shit", "127.0.0.1:9000")

	for _, v := range svrMap.svrMap["shit"].Service {
		if v.ServiceAddr == "127.0.0.1:9000" {
			if !v.IsExpire {
				t.Error("fail")
			}
		}
	}
}

func BenchmarkServiceMap_DelMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		svrMap.AddMap("i", strconv.Itoa(i))
		svrMap.DelMap("i", strconv.Itoa(i))
	}
}

func testA() {
	//svrMap.AddMap("shit","127.0.0.1:9000")
	//fmt.Println(svrMap)
	//svrMap.AddMap("shit2","127.0.0.1:9002")
	//fmt.Println(svrMap)
	//svrMap.DelMap("shit2","127.0.0.1:9002")
}
