package example

import (
	"context"
	"fmt"
	"sandwich"
	"time"
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
	start := time.Now()
	defer func() {
		cost := time.Now().UnixMilli() - start.UnixMilli()
		duration := time.Duration(cost) * time.Millisecond
		fmt.Println("cost:" + duration.String())
	}()
	fmt.Println("start:" + start.String())
	fmt.Println("WrapperMethod")
	for _, param := range invoke.Params {
		if param.Name == "key" {
			if v, ok := param.Value.(string); ok {
				param.Value = "prefix:" + v
			}
		}
		if param.Name == "keys" {
			if v, ok := param.Value.([]string); ok {
				for i := range v {
					v[i] = "prefix:" + v[i]
				}
			}
		}
	}
	// 需要统计耗时 所以这里必须显示调用
	invoke.Process()
}

// IWrapperOp 定义要生成包装实现的方法
type IWrapperOp interface {
	// todo 可以生成扩展的接口方法
	//IOrigin
	//fmt.Stringer
	DoSomething()
	DoSomething1(ctx context.Context, key, val string) error
	DoSomething2(ctx context.Context, keys ...string) (string, error)
}
