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
	fmt.Printf("DoSomething1, key->  %s \n", key)
	return nil
}

func (o *Origin) DoSomething2(ctx context.Context, keys ...string) (string, error) {
	fmt.Printf("DoSomething2, keys->  %+v \n", keys)
	return "", nil
}
