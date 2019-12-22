package model

//http://gorm.book.jasperxu.com/
import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"grpc-lb/internal/common/tool"
	_ "grpc-lb/internal/common/tool"
)

var templateStoreShowCacheKey = "template_store_show_cache"

type TemplateStoreModel struct {
	TemplateId       string
	DesignerId       int32
	Price            int32
	Version          int32
	Star             int32
	Level            int32
	Sort             int32
	IsVipFree        int32
	TemplateInfo     string
	TemplateProperty string
}

func NewTemplateStoreModel() *TemplateStoreModel {
	return &TemplateStoreModel{}
}

func (m *TemplateStoreModel) TableName() string {
	return "platv5_template_store"
}

func (m *TemplateStoreModel) Get(templateId string) error {
	// 如果有缓存，则先读取缓存数据
	conn := tool.RedisPool.Get()
	defer conn.Close()
	resp, err := redis.Bytes(conn.Do("get", templateStoreShowCacheKey))
	if err == nil {
		return json.Unmarshal(resp, m)
	} else {
		fmt.Println(err)
	}

	tool.Mysql.First(m, "template_id = ? ", templateId)
	if m.TemplateId == "" {
		return errors.New("empty rows")
	}
	mJson, _ := json.Marshal(m)
	_, _ = conn.Do("set", templateStoreShowCacheKey, mJson)

	return nil
}
