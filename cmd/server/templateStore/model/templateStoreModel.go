package model
//http://gorm.book.jasperxu.com/
import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type TemplateStoreModel struct {
	TemplateId           string
	DesignerId           int32
	Price                int32
	Version              int32
	Star                 int32
	Level                int32
	Sort                 int32
	IsVipFree            int32
	TemplateInfo         string
	TemplateProperty     string
}

func (m *TemplateStoreModel) Get(templateId string) error {
	db, err := gorm.Open("mysql","root:1234@/makaplatv4?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.First(m, "template_id = ? ",templateId)

	if m.TemplateId == "" {
		return errors.New("empty rows")
	}

	return nil
}

