package Vcache

import (
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	//因为GeterFunc函数类型实现了Getter里定义的Get方法
	//所以可以在Get方法里回调GeterFunc函数
	//将一个匿名函数类型转换为Getter接口，并调用接口的方法
	//其实是在调用匿名回调函数
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	except := []byte("key")
	//回调函数是将字符串转换为切片
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, except) {
		t.Errorf("callback failed")
	}
}
