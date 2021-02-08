package main

import (
	"reflect"
	"strings"
)

/**
快速设置结构体属性字段的值，设置成功返回true，失败返回false
stu：传入结构体指针
field：模糊字段名称，可以是属性字段全称或部分字符串，也可以是属性字段json标签值的全值或部分字符串
val：用作替换的值
*/
func StuSetFieldVal(stu interface{}, field string, val interface{}) bool {
	typ := reflect.TypeOf(stu).Elem()
	for i := 0; i < typ.NumField(); i++ {
		if strings.Index(typ.Field(i).Name, field) != -1 || strings.Index(typ.Field(i).Tag.Get("json"), field) != -1 {
			vl := reflect.ValueOf(stu).Elem()
			vl.FieldByName(typ.Field(i).Name).Set(reflect.ValueOf(val))
			return true
		}
	}
	return false
}

/**
快速获取结构体属性字段的值
传参：
	stu：传入结构体指针
	field：模糊字段名称，可以是属性字段全称或部分字符串，也可以是属性字段json标签值的全值或部分字符串
返回：
	val：获取到的值，拿到该值后需val.(原类型)转换为原类型使用
	success：获取成功为true，失败为false
*/
func StuGetFieldVal(stu interface{}, field string) (val interface{}, success bool) {
	typ := reflect.TypeOf(stu).Elem()
	for i := 0; i < typ.NumField(); i++ {
		if strings.Index(typ.Field(i).Name, field) != -1 || strings.Index(typ.Field(i).Tag.Get("json"), field) != -1 {
			vl := reflect.ValueOf(stu).Elem()
			val = vl.FieldByName(typ.Field(i).Name).Interface()
			success = true
			return
		}
	}
	success = false
	return
}
