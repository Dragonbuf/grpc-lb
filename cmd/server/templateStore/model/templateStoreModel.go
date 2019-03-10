package model

//http://gorm.book.jasperxu.com/
import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"grpc-lb/tool"
	_ "grpc-lb/tool"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
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
	redisClient := tool.RedisPool.Get()
	defer redisClient.Close()
	resp, err := redis.String(redisClient.Do("get",templateStoreShowCacheKey))
	if err == nil {
		err = json.Unmarshal([]byte(resp),m)
		if err == nil {
			return nil
		}
	}

	defer tool.Mysql.Close()
	tool.Mysql.First(m, "template_id = ? ", templateId)
	if m.TemplateId == "" {
		return errors.New("empty rows")
	}

	return nil
}

func (m *TemplateStoreModel) Update()  {
	tx := tool.Mysql.Begin()
	defer tool.Mysql.Close()
	tx.Rollback()
}
