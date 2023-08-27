package vm

import (
	"github.com/vvbae/monkey/code"
	"github.com/vvbae/monkey/object"
)

type Frame struct {
	cl          *object.Closure // points to the compiled function referenced by the frame inside the closure
	ip          int             // in the frame for the fn
	basePointer int             // stack pointer’s value before we execute a function
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	f := &Frame{cl: cl, ip: -1, basePointer: basePointer}
	return f
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
