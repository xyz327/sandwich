package sandwich

import (
	"context"
	sync "sync"
)

type Wrapper[T any] interface {
	// GetDelegate 获取包装的原始对象
	GetDelegate() T
	// WrapperMethod 包装方法的回调，这里可以做一些处理
	WrapperMethod(ctx context.Context, invoke *Invoke)
}

type Valued struct {
	Name  string
	Value any
}

type Invoke struct {
	once sync.Once
	Ctx  context.Context
	// 方法名
	MethodName string
	// 方法入参
	Params []*Valued
	// 方法返回值
	Returns    []*Valued
	_DoProcess func() []*Valued
	_proceed   bool
}

func (p *Invoke) Process() []*Valued {
	p._proceed = true
	return p._DoProcess()
}
func (p *Invoke) SetProcess(fn func() []*Valued) {
	p.once.Do(func() {
		p._DoProcess = fn
	})
}

func (p *Invoke) IsProceed() bool {
	return p._proceed
}
