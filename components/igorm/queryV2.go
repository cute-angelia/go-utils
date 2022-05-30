package igorm

import (
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type Query struct {
	subset map[string]subset
}

// subset 为字段查询结构体
type subset struct {
	// Op MySQL中查询条件，如like,=,in
	op string
	// Val 查询条件对应的值
	val interface{} // 值
}

func NewQuery() *Query {
	return &Query{
		subset: make(map[string]subset),
	}
}

// Where 默认不查询0值
func (q *Query) Where(key string, op string, val interface{}) *Query {
	if !reflect.ValueOf(val).IsZero() {
		q.subset[key] = subset{
			op:  op,
			val: val,
		}
	}
	return q
}

// WhereCanZero 可以查询0值
func (q *Query) WhereCanZero(key string, op string, val interface{}) *Query {
	q.subset[key] = subset{
		op:  op,
		val: val,
	}
	return q
}

// BuildQuery 构建sql和绑定的参数
func (q *Query) BuildQuery() (sql string, binds []interface{}) {
	sql = "1=1"
	binds = make([]interface{}, 0, len(q.subset))
	for field, cond := range q.subset {
		// 说明有表的数据
		if strings.Contains(field, ".") {
			arr := strings.Split(field, ".")
			if len(arr) != 2 {
				return
			}
			field = "`" + arr[0] + "`.`" + arr[1] + "`"
		} else {
			field = "`" + field + "`"
		}
		switch strings.ToLower(cond.op) {
		case "like":
			if cond.val != "" {
				sql += " AND " + field + " like ?"
				cond.val = "%" + cond.val.(string) + "%"
			}
		case "%like":
			if cond.val != "" {
				sql += " AND " + field + " like ?"
				cond.val = "%" + cond.val.(string)
			}
		case "like%":
			if cond.val != "" {
				sql += " AND " + field + " like ?"
				cond.val = cond.val.(string) + "%"
			}
		case "in", "not in":
			sql += " AND " + field + cond.op + " (?) "
		case "between":
			sql += " AND " + field + cond.op + " ? AND ?"
			val := cast.ToStringSlice(cond.val)
			binds = append(binds, val[0], val[1])
			continue
		case "exp":
			sql += " AND " + field + " ? "
			cond.val = gorm.Expr(cond.val.(string))
		default:
			sql += " AND " + field + cond.op + " ? "
		}
		binds = append(binds, cond.val)
	}
	return
}
