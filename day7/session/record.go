package session

import (
	"day7/clause"
	"day7/log"
	"reflect"
)

// 操作记录

// Insert(&u1)
func (s *Session) Insert(model interface{}) (int64, error) {
	table := s.Model(model).cacheSchema
	var tempVars []interface{}
	modelValue := reflect.Indirect(reflect.ValueOf(model))
	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		tempVars = append(tempVars, field.Interface())
	}
	s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
	s.clause.Set(clause.VALUES, tempVars...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Filter([]&u1)
func (s *Session) Filter(result interface{}) error {
	// 查询需要: 表名 字段
	// 取出单个元素 type.elem
	modelSlice := reflect.Indirect(reflect.ValueOf(result))
	modelType := modelSlice.Type().Elem()
	// 获取value通过value得到interface
	modelValue := reflect.New(modelType)
	table := s.Model(modelValue.Interface()).cacheSchema
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDER, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		log.Error(err)
		return err
	}

	for rows.Next() {
		// 新建dest
		newModelValue := reflect.New(modelType).Elem()
		var fieldAddrValue []interface{}
		for i := 0; i < newModelValue.NumField(); i++ {
			fieldAddrValue = append(fieldAddrValue, newModelValue.Field(i).Addr().Interface())
		}

		if err := rows.Scan(fieldAddrValue...); err != nil {
			return err
		}

		// 在反射中加入值
		modelSlice.Set(reflect.Append(modelSlice, newModelValue))

	}
	return rows.Close()
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var t []interface{}
	s.clause.Set(clause.WHERE, append(append(t, desc), args...)...)
	return s
}

func (s *Session) LIMIT(n int) *Session {
	s.clause.Set(clause.LIMIT, n)
	return s
}

func (s *Session) Update(kv ...interface{}) (int64, error) {
	m := map[string]interface{}{}
	for i := 0; i < len(kv); i += 2 {
		if i+1 <= len(kv)-1 {
			m[kv[i].(string)] = kv[i+1]
		}
	}

	s.clause.Set(clause.UPDATE, s.cacheSchema.Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		log.Error(err)
	}
	return result.RowsAffected()

}
