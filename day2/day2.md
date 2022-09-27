# Sqlite
**轻量级数据库, 不需要启动独立服务, 数据存储在单一的磁盘文件中
使用场景: 并发小, 数据量级小**

```go
import(
    _ "github.com/mattn/go-sqlite3" // 包导入时会注册sqlite驱动
    "database/sql" // go提供标准库用于和数据库交互
)
```

- **两个对象, 交互用户和交互数据库**
```
交互用户对象(gdb)具有方法: 
    创建数据库连接对象
    创建交互数据库对象session
    断开连接处理
交互数据库对象(session): 
    创建session对象
    保存sql及参数(复用存储sql变量)
    清空sql变量
    封装sqlite原生方法(添加打印sql语句及执行完毕清空sql)
```
- Exec()用于执行SQL语句, 如果是查询语句不会返回相关记录。查询使用Query()和 QueryRow()。



___
### 自定义日志
**作用: 区分颜色 显示报错行数**
实际通过log.New()通过参数, 设置终端输出前缀及颜色, 显示报错行数和时间
```go
var ErrorLog = log.New(os.stdout, "\033[31m[error]\033", log.Lshortfile|log.LstdFlags)
// 导出方法
var Error = ErrorLog.Println
// 使用
log.Error("test")
```
___ 
### 对象表结构映射Schema
- 为了使ORM框架可以兼容更多的数据库, 创建了Dialect接口(实现差异方法)。例如解析go中不同反射类型对应的数据库类型字符串的方法
- 导包时通过init方法自动注册
 ```var _ Dialect = (*sqlite3)(nil) // 静态检测sqlite3是否实现全部Dialect接口方法 ```

- 定义Parse方法, 将结构体转换为数据库结构

 ```go
// 字段Field和数据库结构Schema

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}
 ```


- 解析操比较耗时, 因此在Session结构里缓存已解析的Model, 当传入Model不一致时才会进行解析
- 封装创建、删除、表是否存在等方法, 通过schema里的Fields列表, 拼接sql
- reflect.DeepEqual 深度比较变量(数组切片等)
- 调用 t.Helper() 让报错信息更准确，有助于定位。**返回外层调用者报错的行号**
- reflect.Value. Addr().Interface() 取地址指针
- reflect.SliceOf 创建一个切片
- TestMain是测试入口函数, m.Run()是阻塞等待测试完成的, 可用于做一些测试前置处理
- 可利用反射查找是否有对应的方法(reflect.Value.MethondByName 获取指定方法), 有则调用, 可实现钩子; 也可利用interface实现, 例如:
```go
type IBeforeQuery interface {
      BeforeQuery(s *Session) error
}

type IAfterQuery interface {
      AfterQuery(s *Session) error
}
.....
等等

//然后修改CallMethod
func (s *Session) CallMethod(method string, value interface{}) {
	 ...
     if i, ok := dest.(IBeforQuery); ok == true {
        i. BeforeQuery(s) 
     }
     ...
	return
}
```
- **注意: reflect.Value对象和reflect.Type对象都有Elem方法, 但是一个是指针指向的值, 一个是返回迭代对象里的一个值**
- **reflect.Value 通过.Type()转换成reflect.Type对象, reflect.Type对象通过New()转换成reflect.value指针对象**
#### Clause构造SQL语句
SELECT语句的构成通常是这样的:
```sql
SELECT col1, col2, ...
    FROM table_name
    WHERE [ conditions ]
    GROUP BY col1
    HAVING [ conditions ]
```
通常由多个子句构成

- 多次调用方法构造属性的, 非常适用链式调用 **(返回对象继续调用)**
- defer中的recover可以使进入宕机流程中的 goroutine 恢复过来, 并且捕获到panic的输入值
- mysql事务中如果遇到ddl(创建表等)会进行一次隐式的commit, 并开启新事务, 导致无法回滚这部分
- t.Run() 可以指定执行方法