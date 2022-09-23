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

