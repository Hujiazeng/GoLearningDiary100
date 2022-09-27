package schema

import (
	"day2/dialect"
	"reflect"
)

// Go转换成数据库格式

type Field struct {
	Name, Type, Tag string
}

// 数据库结构格式
type Schema struct {
	Model      interface{}       // 保存go原结构
	Name       string            // 表名
	Fields     []*Field          // 字段列表
	FieldNames []string          // 字段名列表
	fieldMap   map[string]*Field //字段名对应字段
}

func (Schema *Schema) GetField(name string) *Field {
	return Schema.fieldMap[name]
}

// 实现解析结构体函数
func Parse(unknowModel interface{}, d dialect.Dialect) *Schema {
	modelValue := reflect.Indirect(reflect.ValueOf(unknowModel))
	modelType := modelValue.Type()

	schema := &Schema{
		Model:    unknowModel,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelValue.NumField(); i++ {
		f := modelType.Field(i)
		v := modelValue.Field(i)
		// 构造field对象
		field := &Field{
			Name: f.Name,
			Type: d.DataTypeOf(reflect.ValueOf(v.Interface())), // 字段的reflectValue
		}
		if v, ok := f.Tag.Lookup("gorm"); ok {
			field.Tag = v
		}

		// field对象添加列表字典中
		schema.Fields = append(schema.Fields, field)
		schema.FieldNames = append(schema.FieldNames, f.Name)
		schema.fieldMap[f.Name] = field
	}
	return schema
}

// 由于ORM期望调用的方式 &User{Name:"Tom", Age: 18}, 故需要按数据库中列的字段顺序, 从对象中取出对应值
func (s *Schema) RecordValues(model interface{}) (recordValues []interface{}) {
	modelValue := reflect.Indirect(reflect.ValueOf(model))

	for _, field := range s.Fields {
		recordValues = append(recordValues, modelValue.FieldByName(field.Name).Interface())
	}
	return
}
