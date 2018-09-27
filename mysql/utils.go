package mysql

import (
	"encoding/json"
	"context"
	"fmt"
)

// ps: interface 全部传指针
// 创建或者更新
func CreateOrUpdate(ctx context.Context, dbName string, table string, data interface{}, id int32) (interface{}, error) {
	orm := GetGormCtx(ctx, dbName)
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
func GetPageData(ctx context.Context, dbName string, table string, page int32, prepage int32, data interface{}, count *int32) {
	orm := GetGormCtx(ctx, dbName)
	offset := (page - 1) * prepage
	orm.Table(table).Limit(prepage).Offset(offset).Find(data)
	orm.Table(table).Count(count)
}

// 转化数据
func Convert(src interface{}, dest interface{}) {
	temp, _ := json.Marshal(src)
	json.Unmarshal(temp, dest)
}
