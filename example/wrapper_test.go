package example

import (
	"context"
	"fmt"
	"testing"
)

func TestGen(t *testing.T) {
	w := &WrapperTest{&Origin{}}

	// 会先执行 WrapperTest 定义的方法输出 WrapperMethod
	// 再执行 Origin 的方法输出 DoSomething
	fmt.Println("DoSomething")
	w.DoSomething()

	fmt.Println("DoSomething1")
	w.DoSomething1(context.Background(), "key", "val")

	fmt.Println("DoSomething2")
	w.DoSomething2(context.Background(), "key1", "key2")
}
