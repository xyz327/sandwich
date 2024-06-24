package example

import (
	"context"
	"fmt"
	"sandwich"
)

var _ sandwich.Wrapper[IOrigin] = (*WrapperTest)(nil)

//go:generate go run ../cmd/gen.go main -input=wrapper.go -type=WrapperTest -output=wrapper_gen.go
type WrapperTest struct {
	origin Origin
}

func (w *WrapperTest) GetDelegate() IOrigin {
	return &w.origin
}

func (w *WrapperTest) WrapperMethod(ctx context.Context, invoke *sandwich.Invoke) {
	fmt.Println("WrapperMethod")
}

type IWrapperOp interface {
	//IOrigin
	//fmt.Stringer
	DoSomething()
	DoSomething1(ctx context.Context, key, val string) error
	DoSomething2(ctx context.Context, key ...string) (string, error)
}
