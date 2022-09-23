package main

import (
	"fmt"
	"reflect"
)

type myfloat64 float64

// 结构体
type unknow struct {
	s1, S2, S3 string
}

func (u unknow) Stringd() string {
	return u.s1 + "-" + u.S2 + "-" + u.S3
}

type T struct {
	A int
	B string
}

func main() {
	var x myfloat64 = 3.14
	fmt.Println(reflect.ValueOf(x))
	v := reflect.ValueOf(x) // 值传递故v修改不到x
	fmt.Println("value:", v)
	fmt.Println("type:", v.Type())
	fmt.Println("kind:", v.Kind())

	fmt.Println(".Float:", v.Float())
	// fmt.Println(".int:", v.Int()) float类型使用会报错
	fmt.Println(".Interface:", v.Interface()) // 从对象返回接口值

	y := v.Interface().(myfloat64)
	fmt.Printf("%T\n", y)

	// 通过反射修改(设置)值
	// v.SetFloat(666) // error: using unaddressable value
	fmt.Println("canset?", v.CanSet()) //false

	// 当v:=relect.ValueOf(x) 函数通过传递一个x值拷贝创建了一个v, 那么v的改变并不能更改原始的x。
	// 想要v的更改能做到用x, 那就必须传递x的地址v=reflect.ValueOf(&x)

	v = reflect.ValueOf(&x)
	fmt.Println("canset?", v.CanSet())                //false
	fmt.Printf("ValueOf传递地址返回反射对象类型: %v\n", v.Type()) // *main.myfloat64
	// 指针仍然是不可设置的, 想要其可设置需要使用Elem()函数, 间接使用指针指向的变量
	v.Elem().SetFloat(666.123) // success
	fmt.Println(v.Elem().Interface())

	// 反射结构
	fmt.Println("--------Struct-------")
	var secretStruct interface{} = unknow{"Xiaoming", "Xiaohu", "xiaoxia"} // 空接口存放了变量指针和类型
	value := reflect.ValueOf(secretStruct)
	typ := reflect.TypeOf(secretStruct)
	knd := value.Kind()
	fmt.Println("value:", value) // 如果结构体有定义string方法, Println会直接调用返回
	fmt.Println("type:", typ)
	fmt.Println("kind: ", knd)

	// 遍历字段
	for i := 0; i < value.NumField(); i++ {
		fmt.Printf("Field Index: %d, Value: %v\n", i, value.Field(i))
	}

	// 调用结构体方法
	result := value.Method(0).Call(nil)
	fmt.Println(result)
	// 尝试修改值
	// value.Field(0).SetString("123") // panic: using unexported field, 只有首字母大写的字段才可被修改
	fmt.Println("CanSet?", value.Field(1).CanSet()) // false
	// 传入地址
	newValue := reflect.ValueOf(&unknow{"xiaomei", "xiaohu", "xiaojia"})
	// 使用Elem指向指针的值
	newValue.Elem().Field(1).SetString("123")
	// 调用方法查看是否成功
	fmt.Println(newValue.Method(0).Call(nil))

	fmt.Println(value.Type() == typ)
	fmt.Println(typ.Name())
}
