package schema

import (
	"go/ast"
	"reflect"
	"vgo/orm/day2-reflect-schema/dialect"
)

/**
Dialect实现了一些特定的SQL语句的转换，接下来我们将要实现ORM框架中最为核心的转换--对象(object)和表(table)的转换。
给定一个任意的对象，转换为关系型数据库中的表结构。
*/

// Field represents a column of database
type Field struct {
	Name string // 字段名
	Type string // 字段类型
	Tag  string // 约束条件
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}       // 被映射的对象
	Name       string            // 表名
	Fields     []*Field          // 字段
	FieldNames []string          // 包含所有的字段名(列名)
	fieldMap   map[string]*Field // 记录字段名和Field的映射关系，方便之后直接使用，无需遍历Fields
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

/**
将任意的对象解析为Schema实例
TypeOf() 和 ValueOf()是reflect包最为基本也是最重要的2个方法，分别用来返回入参的类型和值。
因为设计的入参是一个对象的指针，因此需要reflect.Indirect()获取指针指向的实例。
modelType.Name() 获取到结构体的名称作为表名
NumField() 获取到实例的字段的个数，然后通过下标获取到特定字段p := modelType.Field(i)
p.Name即字段名，p.Type即字段类型，通过(Dialect).DataTypeOf()转换为数据库的字段类型，p.Tag即额外的约束条件
*/
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("vgoorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
