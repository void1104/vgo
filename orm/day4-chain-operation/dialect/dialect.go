package dialect

import "reflect"

/**
实现ORM映射的第一步，需要思考如何将GO语言的类型映射为数据库中的类型

不同数据库支持的数据类型也是有差异的，即使功能相同，在SQL语句的表达上也可能有差异。
ORM框架往往需要兼容多种数据库，因此我们需要将差异的一部分提取出来，每一种数据库分别实现，实现最大程度的复用和解耦。
这部分代码称之为dialect
*/
var dialectsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
