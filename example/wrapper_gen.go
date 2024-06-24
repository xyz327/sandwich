// Code generated by go generate; DO NOT EDIT.

package example

import (
	"context"
	"sandwich"
)

func (c *WrapperTest) _delegateCall(invoke *sandwich.Invoke) {
	defer func() {
		if !invoke.IsProceed() {
			invoke.Process()
		}
	}()
	c.WrapperMethod(invoke.Ctx, invoke)
}

// implements IWrapperOp

func (c *WrapperTest) DoSomething() {
	_params := []*sandwich.Valued{}
	invoke := &sandwich.Invoke{Ctx: context.Background(), MethodName: "DoSomething", Params: _params}
	invoke.SetProcess(func() []*sandwich.Valued {
		c.GetDelegate().DoSomething()
		return invoke.Returns

	})
	c._delegateCall(invoke)
	return
}

func (c *WrapperTest) DoSomething1(ctx context.Context, key string, val string) (_error error) {
	_params := []*sandwich.Valued{
		{Name: "ctx", Value: ctx}, {Name: "key", Value: key}, {Name: "val", Value: val},
	}
	invoke := &sandwich.Invoke{Ctx: ctx, MethodName: "DoSomething1", Params: _params}
	invoke.SetProcess(func() []*sandwich.Valued {

		_error = c.GetDelegate().DoSomething1(_params[0].Value.(context.Context), _params[1].Value.(string), _params[2].Value.(string))
		invoke.Returns = make([]*sandwich.Valued, 1)
		invoke.Returns[0] = &sandwich.Valued{Name: "_error", Value: _error}

		return invoke.Returns

	})
	c._delegateCall(invoke)
	return
}

func (c *WrapperTest) DoSomething2(ctx context.Context, key ...string) (_string string, _error error) {
	_params := []*sandwich.Valued{
		{Name: "ctx", Value: ctx}, {Name: "key", Value: key},
	}
	invoke := &sandwich.Invoke{Ctx: ctx, MethodName: "DoSomething2", Params: _params}
	invoke.SetProcess(func() []*sandwich.Valued {

		_string, _error = c.GetDelegate().DoSomething2(_params[0].Value.(context.Context), _params[1].Value.([]string)...)
		invoke.Returns = make([]*sandwich.Valued, 2)
		invoke.Returns[0] = &sandwich.Valued{Name: "_string", Value: _string}
		invoke.Returns[1] = &sandwich.Valued{Name: "_error", Value: _error}

		return invoke.Returns

	})
	c._delegateCall(invoke)
	return
}