package dialect

import (
	"fmt"
	"reflect"
	"time"
)

/**
对sqlite3的支持
sqlite3 实现虽然比较繁琐，但是整体逻辑还是非常清晰的。DataTypeOf将GO语言映射为SQLite的
数据类型。TableExistSQL返回了在SQLite中判断表tableName是否存在的SQL语句。
实现了 init 函数，包在第一次加载时，会将 sqlite3 的 dialect 自动注册到全局
*/
type sqlite3 struct{}

var _ Dialect = (*sqlite3)(nil)

func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' ans name = ?", args
}
