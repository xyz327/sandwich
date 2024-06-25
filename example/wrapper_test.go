package example

import "testing"

func TestGen(t *testing.T) {
	w := &WrapperTest{&Origin{}}

	// 会先执行 WrapperTest 定义的方法输出 WrapperMethod
	// 再执行 Origin 的方法输出 DoSomething
	w.DoSomething()
}
