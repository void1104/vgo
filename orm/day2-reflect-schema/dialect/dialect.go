package dialect

import "reflect"

/**
实现ORM映射的第一步，需要思考如何将GO语言的类型映射为数据库中的类型

不同数据库支持的数据类型也是有差异的，即使功能相同，在SQL语句的表达上也可能有差异。
ORM框架往往需要兼容多种数据库，因此我们需要将差异的一部分提取出来，每一种数据库分别实现，实现最大程度的复用和解耦。
这部分代码称之为dialect
*/

var dialectsMap = map[string]Dialect{}

/**
Dialect接口包含2个方法：
1. DataTypeOf 用于将GO语言的类型转换为该数据库的数据模型
2. TableExistSQL 返回某个表是否存在的SQL语句，参数是表名(table)
当然，不同数据库之间的差异远远不止这两个地方，随着ORM框架功能的增多，dialect的实现也会逐渐丰富起来
同时框架的其他部分也不会受到影响。
同时，声明了 RegisterDialect 和 GetDialect 两个方法用于注册和获取dialect实例。如果新增加对某个数据库的支持
那么调用 RegisterDialect 即可注册到全局
*/
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
