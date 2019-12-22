package etcd_web_watch

import (
	"encoding/json"
	"fmt"
	_ "github.com/coreos/etcd/mvcc/mvccpb"
	_ "net/http"
	"strconv"
	_ "strings"
	"testing"
)

// 可以起本地 8000 端口，查看服务列表
func TestStart(t *testing.T) {
	Start()
}

func TestServiceMap_AddMap(t *testing.T) {
	AddMap("shit", "127.0.0.1:9000")

	if TotalService != 1 {
		t.Error("测试失败")
	}

	AddMap("shit", "127.0.0.1:9001")
	if TotalService != 2 {
		t.Error("测试失败")
	}
}

func BenchmarkServiceMap_AddMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AddMap("i", strconv.Itoa(i))
	}
}

func BenchmarkIndexHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := i
		a++
	}
}

func TestServiceMap_DelMap(t *testing.T) {
	AddMap("shit", "127.0.0.1:9000")
	DelMap("shit", "127.0.0.1:9000")

	jsonD, _ := json.Marshal(svrMap)
	fmt.Println(string(jsonD))
	for _, v := range Service {
		if ServiceAddr == "127.0.0.1:9000" {
			t.Error("fail")
		}
	}
}

func BenchmarkServiceMap_DelMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AddMap("i", strconv.Itoa(i))
		DelMap("i", strconv.Itoa(i))
	}
}

func testA() {
	//svrMap.AddMap("shit","127.0.0.1:9000")
	//fmt.Println(svrMap)
	//svrMap.AddMap("shit2","127.0.0.1:9002")
	//fmt.Println(svrMap)
	//svrMap.DelMap("shit2","127.0.0.1:9002")
}
