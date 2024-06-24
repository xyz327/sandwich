package sandwich

import (
	"context"
	sync "sync"
)

type Wrapper[T any] interface {
	GetDelegate() T
	WrapperMethod(ctx context.Context, invoke *Invoke)
}

type Valued struct {
	Name  string
	Value any
}

type Invoke struct {
	once       sync.Once
	MethodName string
	Ctx        context.Context
	Params     []*Valued
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
