package igorm

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

// ps: interface 全部传指针
// 创建或者更新
func CreateOrUpdate(orm *gorm.DB, table string, data map[string]interface{}, id int32) (interface{}, error) {
	if id > 0 {
		if orm.Table(table).Where("id= ?", id).Updates(data).RowsAffected == 0 {
			return nil, fmt.Errorf("更新错误")
		}
	} else {
		orm.Table(table).Create(data)
	}
	return data, nil
}

// 查询分页数据
func GetPageData(orm *gorm.DB, table string, page int, prepage int, model interface{}, count int64) (interface{}, int64) {
	offset := (page - 1) * prepage
	orm.Table(table).Limit(prepage).Offset(offset).Find(&model)
	orm.Table(table).Count(&count)
	return model, count
}

// 转化数据 dest => &dest
func Convert(src interface{}, dest interface{}) {
	temp, _ := json.Marshal(src)
	json.Unmarshal(temp, dest)
}

// gorm updates 对 model 为 0 的数据不处理， 这里转化为 map 对象处理
func ConvertMap(in interface{}, noKey []string) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(in)
	json.Unmarshal(inrec, &inInterface)

	if _, ok := inInterface["id"]; ok {
		delete(inInterface, "id")
	}

	if _, ok := inInterface["uid"]; ok {
		delete(inInterface, "uid")
	}

	for _, i2 := range noKey {
		if _, ok := inInterface[i2]; ok {
			delete(inInterface, i2)
		}
	}

	for k, v := range inInterface {
		if fmt.Sprintf("%v", v) == "" {
			delete(inInterface, k)
		}
	}

	return inInterface
}

func QueryGenerate(orm *gorm.DB, key string, opt string, value interface{}) *gorm.DB {
	switch v := value.(type) {
	case string:
		if len(v) > 0 {
			if opt == "like" {
				orm = orm.Where(fmt.Sprintf("%s like ?", key), "%"+v+"%")
			} else {
				orm = orm.Where(fmt.Sprintf("%s %s ?", key, opt), v)
			}
		}
	case int:
		if v > 0 {
			orm = orm.Where(fmt.Sprintf("%s %s ?", key, opt), v)
		}
	case int32:
		if v > 0 {
			orm = orm.Where(fmt.Sprintf("%s %s ?", key, opt), v)
		}
	case []int32:
		if len(v) > 0 {
			orm = orm.Where(fmt.Sprintf("%s %s (?)", key, opt), v)
		}
	default:
		orm = orm.Where(fmt.Sprintf("%s %s ?", key, opt), v)
	}

	return orm
}
