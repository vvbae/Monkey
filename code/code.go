package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	// arithmetic
	OpSub
	OpMul
	OpDiv
	// boolean
	OpTrue
	OpFalse
	// logical
	OpEqual
	OpNotEqual
	OpGreaterThan
	// prefix
	OpMinus
	OpBang
	// jump
	OpJumpNotTruthy
	OpJump
	// null
	OpNull
	// global bindings
	OpGetGlobal
	OpSetGlobal
	// data structures
	OpArray
	OpHash
	// index
	OpIndex
	// function calls
	OpCall
	// return
	OpReturnValue
	OpReturn
	// local bindings
	OpGetLocal
	OpSetLocal
	// builtin scope
	OpGetBuiltin
	// closure
	OpClosure
	// free variables
	OpGetFree
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},      // uint16, ~65535
	OpAdd:           {"OpAdd", []int{}},            // +
	OpPop:           {"OpPop", []int{}},            // pop value from the stack
	OpSub:           {"OpSub", []int{}},            // -
	OpMul:           {"OpMul", []int{}},            // *
	OpDiv:           {"OpDiv", []int{}},            // /
	OpTrue:          {"OpTrue", []int{}},           // true
	OpFalse:         {"OpFalse", []int{}},          // false
	OpEqual:         {"OpEqual", []int{}},          // =
	OpNotEqual:      {"OpNotEqual", []int{}},       // !=
	OpGreaterThan:   {"OpGreaterThan", []int{}},    // >
	OpMinus:         {"OpMinus", []int{}},          // negative
	OpBang:          {"OpBang", []int{}},           // !
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}}, // jump to location if not true
	OpJump:          {"OpJump", []int{2}},          // jump to location
	OpNull:          {"OpNull", []int{}},           // null
	OpGetGlobal:     {"OpGetGlobal", []int{2}},     // bind value to global variable
	OpSetGlobal:     {"OpSetGlobal", []int{2}},     // get value of the global variable
	OpArray:         {"OpArray", []int{2}},         // number of elements ~65535
	OpHash:          {"OpHash", []int{2}},          // length of the hashmap
	OpIndex:         {"OpIndex", []int{}},          // index and array are stored before
	OpCall:          {"OpCall", []int{1}},          // function call with arg number
	OpReturnValue:   {"OpReturnValue", []int{}},    // return the value sitting on top of the stack
	OpReturn:        {"OpReturn", []int{}},         // just return, no value
	OpGetLocal:      {"OpGetLocal", []int{1}},      // bind value to local variable
	OpSetLocal:      {"OpSetLocal", []int{1}},      // get value of local variable
	OpGetBuiltin:    {"OpGetBuiltin", []int{1}},    // index of the function
	OpClosure:       {"OpClosure", []int{2, 1}},    // constant index: for locating *object.CompiledFunction, free variable count
	OpGetFree:       {"OpGetFree", []int{1}},       // number of free variables
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

// encodes the operands of a bytecode instruction
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}

	return instruction
}

// decode operands from a bytecode instruction
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}
