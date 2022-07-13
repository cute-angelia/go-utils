package igorm

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

// CreateOrUpdate 创建或者更新
// ps: interface 全部传指针
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

// GetPageData 查询分页数据
func GetPageData(orm *gorm.DB, tableName string, page int, prepage int, models interface{}) (interface{}, int64) {
	count := int64(0)
	offset := (page - 1) * prepage
	orm.Table(tableName).Count(&count)
	orm.Table(tableName).Limit(prepage).Offset(offset).Find(&models)
	return models, count
}

// GetPageDataV2 查询分页内容
func GetPageDataV2(orm *gorm.DB, tableName string, conds Conds, page int, perPage int) (interface{}, int64) {
	total := int64(0)
	var models interface{}
	sql, binds := BuildQuery(conds)
	db := orm.Table(tableName).Where(sql, binds...)
	db.Count(&total)
	db.Order("id desc").Offset((page - 1) * perPage).Limit(perPage).Find(&models)
	return models, total
}

// GetPageDataV3 查询分页内容
// 例子：总榜数据
// sql, binds := igorm.NewQuery().Where("type", "=", itype).BuildQuery()
// db := orm.Table(userScore.TableName()).Order("score desc").Where(sql, binds...)
//
// ToplistAllModel := []wwxcUserScore.WwxcUserScoreModel{}
// list, _ := igorm.GetPageDataV3(db, ToplistAllModel, page, perpage)
func GetPageDataV3(db *gorm.DB, model interface{}, page int, perPage int) (interface{}, int64) {
	total := int64(0)
	db.Count(&total)
	db.Offset((page - 1) * perPage).Limit(perPage).Find(&model)
	return model, total
}

// Convert 转化数据 dest => &dest
func Convert(src interface{}, dest interface{}) {
	temp, _ := json.Marshal(src)
	json.Unmarshal(temp, dest)
}

// ConvertMap gorm updates 对 model 为 0 的数据不处理， 这里转化为 map 对象处理
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
