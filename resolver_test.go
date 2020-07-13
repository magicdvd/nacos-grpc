package grpc

import (
	"testing"
)

func Test_newNacosResolver(t *testing.T) {
	// _, err := grpc.Dial("http://aaa.bbb.com/asdfsadf?asdfasdf", grpc.WithInsecure())
	// if err != nil {
	// 	t.Error(err)
	// }
	t.Log(Target("http://nacos:nacos@127.0.0.1:8848/nacos", "hello"))
}
