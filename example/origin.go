package example

import (
	"context"
	"fmt"
)

var _ IOrigin = (*Origin)(nil)

type IOrigin interface {
	DoSomething()
	DoSomething1(ctx context.Context, key, val string) error
	DoSomething2(ctx context.Context, key ...string) (string, error)
}
type Origin struct {
}

func (o *Origin) DoSomething() {
	fmt.Println("DoSomething")
}

func (o *Origin) DoSomething1(ctx context.Context, key, val string) error {
	//TODO implement me
	panic("implement me")
}

func (o *Origin) DoSomething2(ctx context.Context, key ...string) (string, error) {
	//TODO implement me
	panic("implement me")
}
