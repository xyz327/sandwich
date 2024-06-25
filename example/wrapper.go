package example

import (
	"context"
	"fmt"
	"sandwich"
)

var _ sandwich.Wrapper[IOrigin] = (*WrapperTest)(nil)

// WrapperTest 定义包装对象  必须实现 sandwich.Wrapper 接口
//
//go:generate go run ../cmd/gen.go main -input=wrapper.go -type=WrapperTest -output=wrapper_gen.go
type WrapperTest struct {
	origin *Origin
}

func (w *WrapperTest) GetDelegate() IOrigin {
	return w.origin
}

// WrapperMethod 包装方法的回调，这里可以做一些处理,目前不会阻塞原流程，
func (w *WrapperTest) WrapperMethod(ctx context.Context, invoke *sandwich.Invoke) {
	fmt.Println("WrapperMethod")
}

// IWrapperOp 定义要生成包装实现的方法
type IWrapperOp interface {
	// todo 可以生成扩展的接口方法
	//IOrigin
	//fmt.Stringer
	DoSomething()
	DoSomething1(ctx context.Context, key, val string) error
	DoSomething2(ctx context.Context, key ...string) (string, error)
}
