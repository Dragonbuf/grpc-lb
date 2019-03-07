package model

//http://gorm.book.jasperxu.com/
import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"grpc-lb/tool"
	_ "grpc-lb/tool"
)

var templateStoreShowCache = "template_store_show_cache"

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

func (m *TemplateStoreModel) TableName() string {
	return "platv5_template_store"
}

func (m *TemplateStoreModel) Get(templateId string) error {

	defer tool.Mysql.Close()

	tool.Mysql.First(m, "template_id = ? ", templateId)

	if m.TemplateId == "" {
		return errors.New("empty rows")
	}

	_, err := tool.RedisPool.Get().Do(templateStoreShowCache)
	if err != nil {

	}

	return nil
}
