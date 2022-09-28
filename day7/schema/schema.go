package schema

import (
	"day7/log"
	"go/ast"
	"reflect"
	"strings"
)

// 数据库结构
type Schema struct {
	Model      interface{} // 结构体模型
	Name       string      // 表名
	Fields     []*Field    // 字段列表
	FieldNames []string    // 字段名列表
}

// 字段
type Field struct {
	Name    string // 字段名
	Type    string // 类型
	Tag     string //标签
	Primary bool   // 主键
	Unique  bool   //	唯一
	NotNull string // 不为NULl
	Comment string // 注释
	UnSign  string // 无符号
	Default string // 默认值
}

type Hero struct {
	NickName string `korm:"primarykey;  type:varchar(99)"`
	Age      int    `korm:"not null; comment:年龄; type:int(12); Default:520; unique;"`
}

// 解析结构体
func Parse(model interface{}) *Schema {
	// 反射遍历字段 构造Field 生成Schema
	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()
	schema := &Schema{}
	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)
		// 过滤
		if !f.Anonymous && ast.IsExported(f.Name) {
			// 解析Tag类型
			v, ok := f.Tag.Lookup("korm")
			if !ok {
				log.Errorf("Struct: %s Field %s, No define korm Tag", modelType.Name(), f.Name)
				return nil
			}
			setting := GetTagSetting(v)
			field := &Field{
				Name:    f.Name,
				Primary: CheckTrue(setting["PRIMARYKEY"], setting["PRIMARY_KEY"]), // 检测字典是否有PRIMARYKEY的key
				NotNull: setting["NOTNULL"],
				Comment: setting["COMMENT"],
				Default: setting["DEFAULT"],
				Type:    setting["TYPE"],
				Unique:  CheckTrue(setting["UNIQUE"]),
				UnSign:  setting["UNSIGN"],
			}

			// 构造schema
			schema.Name = modelType.Name()
			schema.FieldNames = append(schema.FieldNames, field.Name)
			schema.Fields = append(schema.Fields, field)
		}
	}
	schema.Model = model
	return schema
}

// 检测是否为True
func CheckTrue(vars ...string) bool {
	n := len(vars)
	for i := 0; i < n; i++ {
		if vars[i] != "" && !strings.EqualFold(vars[i], "false") {
			return true
		}
	}
	return false
}

// 解析Tag
func GetTagSetting(t string) map[string]string {
	setting := map[string]string{}
	sL := strings.Split(t, ";")
	for i := 0; i < len(sL); i++ {
		kvL := strings.Split(sL[i], ":")
		key := strings.TrimSpace(strings.ToUpper(kvL[0]))
		if len(kvL) > 1 {
			setting[key] = kvL[1]
		} else {
			setting[key] = key
		}
	}

	// 未指定类型报错
	if _, ok := setting["TYPE"]; !ok {
		log.Error("no type")
	}
	return setting
}
