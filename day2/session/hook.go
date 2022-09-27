package session

// 定义hook接口类

type IAfterFind interface {
	AfterFind(s *Session) error
}

func (s *Session) CallMethod(model interface{}) {
	// 反射获取model
	if f, ok := model.(IAfterFind); ok {
		f.AfterFind(s)
	}
}
