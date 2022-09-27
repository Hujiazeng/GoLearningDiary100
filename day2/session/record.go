package session

import (
	"day2/clause"
	"errors"
	"reflect"
)

// 实现记录增删改查相关逻辑

func (s *Session) Insert(models ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)

	for _, model := range models {
		// 解析字段
		cacheTable := s.Model(model).cacheTable
		// 生成sql语句
		s.clause.Set(clause.INSERT, cacheTable.Name, cacheTable.FieldNames)
		// 得到所有插入的字段值列表
		recordValues = append(recordValues, cacheTable.RecordValues(model))
	}

	// 多条插入sql  [[值列表], [值列表]] => "insert into table (fields) values(?,?,?,?), (?,?,?,?)";
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()

}

func (s *Session) Update(m map[string]interface{}) (int64, error) {
	s.clause.Set(clause.UPDATE, s.cacheTable.Name, m)
	sql, vars := s.clause.Build(clause.UPDATE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.cacheTable.Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.cacheTable.Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *Session) Find(emptySlice interface{}) error {
	sliceValue := reflect.Indirect(reflect.ValueOf(emptySlice))
	// 取出单个类型: 只有reflect.Type的Elem取出单个元素, reflect.value的Elem是取指针指向值
	modelType := sliceValue.Type().Elem()
	// 解析获取model的schema
	modelValue := reflect.New(modelType).Elem()
	schema := s.Model(modelValue.Addr().Interface()).cacheTable

	// 执行sql搜索
	s.clause.Set(clause.SELECT, schema.Name, schema.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	defer rows.Close()
	// 遍历每条查询记录
	for rows.Next() {
		values := make([]interface{}, 0) // 暂存字段地址
		// 存储字段地址
		for _, fieldName := range schema.FieldNames {
			values = append(values, modelValue.FieldByName(fieldName).Addr().Interface())
		}
		// scan依次写入地址
		if err = rows.Scan(values...); err != nil {
			return err
		}
		// 字段地址值改变, modelValue发生改变, 将修改后的modelValue存入参数
		sliceValue.Set(reflect.Append(sliceValue, modelValue))

		s.CallMethod(schema.Model)
	}
	return nil
}

// 仅查询一条记录
func (s *Session) First(value interface{}) error {
	// 构建一个slice
	modelValue := reflect.Indirect(reflect.ValueOf(value))
	modelSlice := reflect.New(reflect.SliceOf(modelValue.Type())).Elem()

	if err := s.Find(modelSlice.Addr().Interface()); err != nil {
		return err
	}
	if modelSlice.Len() == 0 {
		return errors.New("no record")
	}
	// 将查询结果写入value
	modelValue.Set(modelSlice.Index(0))
	return nil

}
