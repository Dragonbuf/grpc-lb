package model

//http://gorm.book.jasperxu.com/
import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"grpc-lb/internal/common/cache"
	"grpc-lb/internal/common/db"
	_ "grpc-lb/internal/common/db"
)

var templateStoreShowCacheKey = "_template_"

type TemplateStoreModel struct {
	TemplateId string
}

func NewTemplateStoreModel() *TemplateStoreModel {
	return &TemplateStoreModel{}
}

func (m *TemplateStoreModel) TableName() string {
	return "template"
}

func (m *TemplateStoreModel) Get(templateId string) error {
	// 如果有缓存，则先读取缓存数据
	conn := cache.RedisPool.Get()
	defer conn.Close()
	resp, err := redis.Bytes(conn.Do("get", templateStoreShowCacheKey+templateId))
	if err == nil {
		return json.Unmarshal(resp, m)
	} else {
		fmt.Println(err)
	}

	db.Mysql.First(m, "template_id = ? ", templateId)
	if m.TemplateId == "" {
		return errors.New("empty rows")
	}
	mJson, _ := json.Marshal(m)
	_, _ = conn.Do("set", templateStoreShowCacheKey+templateId, mJson)

	return nil
}
