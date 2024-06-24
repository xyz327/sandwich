package example

import "testing"

func TestGen(t *testing.T) {
	w := &WrapperTest{Origin{}}

	w.DoSomething()
}
