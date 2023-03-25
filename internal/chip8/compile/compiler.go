package compile

import (
	"chip8-emu/internal/chip8/parse"
	"strconv"
)

type Compiler struct {
	Instructions *InstructionSet
	tokens       []parse.Token
}

func NewCompiler(tokens []parse.Token) *Compiler {
	return &Compiler{
		tokens:       tokens,
		Instructions: NewInstructionSet(tokens),
	}
}

func (c Compiler) Compile() []byte {
	opcodes := []byte{}
	for _, inst := range c.Instructions.Instructions {
		switch inst.Tokens[0].Type {
		case parse.CLS:
			opcodes = append(opcodes, []byte{0x00, 0xE0}...)
		case parse.SYSCALL:
			callAddress := valueToInt(inst.Tokens[1])
			op := 0x0000 | (0xFFF & callAddress)
			first := byte(op >> 4)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.CALL:
			callAddress := valueToInt(inst.Tokens[1])
			op := 0x2000 | (0xFFF & callAddress)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.RET:
			opcodes = append(opcodes, []byte{0x00, 0xEE}...)
		case parse.JMP:
			callAddress := valueToInt(inst.Tokens[1])
			op := 0x1000 | (0xFFF & callAddress)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.RJMP:
			callAddress := valueToInt(inst.Tokens[1])
			op := 0xB000 | (0xFFF & callAddress)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.SEQ:
			switch inst.Format {
			case CMD_REG_VAL:
				reg := valueToInt(inst.Tokens[3])
				val := valueToInt(inst.Tokens[6])
				op := 0x3000 | ((reg << 8) & 0x0F00) | (val & 0x00FF)
				first := byte((op & 0xFF00) >> 8)
				second := byte(op & 0xFF)
				opcodes = append(opcodes, []byte{first, second}...)
			case CMD_REG_REG:
				reg1 := valueToInt(inst.Tokens[3])
				reg2 := valueToInt(inst.Tokens[8])
				op := 0x5000 | ((reg1 << 8) & 0x0F00) | ((reg2 << 4) & 0x00F0)
				first := byte((op & 0xFF00) >> 8)
				second := byte(op & 0xFF)
				opcodes = append(opcodes, []byte{first, second}...)
			}
		case parse.SNEQ:
			switch inst.Format {
			case CMD_REG_VAL:
				reg := valueToInt(inst.Tokens[3])
				val := valueToInt(inst.Tokens[6])
				op := 0x4000 | ((reg << 8) & 0x0F00) | (val & 0x00FF)
				first := byte((op & 0xFF00) >> 8)
				second := byte(op & 0xFF)
				opcodes = append(opcodes, []byte{first, second}...)
			case CMD_REG_REG:
				reg1 := valueToInt(inst.Tokens[3])
				reg2 := valueToInt(inst.Tokens[8])
				op := 0x9000 | ((reg1 << 8) & 0x0F00) | ((reg2 << 4) & 0x00F0)
				first := byte((op & 0xFF00) >> 8)
				second := byte(op & 0xFF)
				opcodes = append(opcodes, []byte{first, second}...)
			}
		case parse.JKP:
			reg := valueToInt(inst.Tokens[3])
			op := 0xE09E | ((reg << 8) & 0x0F00)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.JKNP:
			reg := valueToInt(inst.Tokens[3])
			op := 0xE0A1 | ((reg << 8) & 0x0F00)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.WK:
			reg := valueToInt(inst.Tokens[3])
			op := 0xF00A | ((reg << 8) & 0x0F00)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.MOV:
			switch inst.Format {
			case CMD_REG_VAL:
				reg := valueToInt(inst.Tokens[3])
				val := valueToInt(inst.Tokens[6])
				op := 0x6000 | ((reg << 8) & 0x0F00) | (val & 0x00FF)
				first := byte((op & 0xFF00) >> 8)
				second := byte(op & 0xFF)
				opcodes = append(opcodes, []byte{first, second}...)
			case CMD_REG_REG:
				reg1 := valueToInt(inst.Tokens[3])
				reg2 := valueToInt(inst.Tokens[8])
				op := 0x8000 | ((reg1 << 8) & 0x0F00) | ((reg2 << 4) & 0x00F0)
				first := byte((op & 0xFF00) >> 8)
				second := byte(op & 0xFF)
				opcodes = append(opcodes, []byte{first, second}...)
			case CMD_ADP_VAL:
				val := valueToInt(inst.Tokens[3])
				op := 0xA000 | (val & 0x0FFF)
				first := byte((op & 0xFF00) >> 8)
				second := byte(op & 0xFF)
				opcodes = append(opcodes, []byte{first, second}...)
			case CMD_REG_DLY:
				reg := valueToInt(inst.Tokens[3])
				op := 0xF007 | ((reg << 8) & 0x0F00)
				first := byte((op & 0xFF00) >> 8)
				second := byte(op & 0xFF)
				opcodes = append(opcodes, []byte{first, second}...)
			case CMD_DLY_REG:
				if inst.Tokens[1].Type == parse.DELAY {
					reg := valueToInt(inst.Tokens[5])
					op := 0xF015 | ((reg << 8) & 0x0F00)
					first := byte((op & 0xFF00) >> 8)
					second := byte(op & 0xFF)
					opcodes = append(opcodes, []byte{first, second}...)
				} else { // SND_DELAY
					reg := valueToInt(inst.Tokens[5])
					op := 0xF018 | ((reg << 8) & 0x0F00)
					first := byte((op & 0xFF00) >> 8)
					second := byte(op & 0xFF)
					opcodes = append(opcodes, []byte{first, second}...)
				}
			}
		case parse.ADD:
			switch inst.Format {
			case CMD_REG_VAL:
				reg := valueToInt(inst.Tokens[3])
				val := valueToInt(inst.Tokens[6])
				op := 0x7000 | ((reg << 8) & 0x0F00) | (val & 0x00FF)
				first := byte((op & 0xFF00) >> 8)
				second := byte(op & 0xFF)
				opcodes = append(opcodes, []byte{first, second}...)
			case CMD_REG_REG:
				reg1 := valueToInt(inst.Tokens[3])
				reg2 := valueToInt(inst.Tokens[8])
				op := 0x8004 | ((reg1 << 8) & 0x0F00) | ((reg2 << 4) & 0x00F0)
				first := byte((op & 0xFF00) >> 8)
				second := byte(op & 0xFF)
				opcodes = append(opcodes, []byte{first, second}...)
			case CMD_ADP_REG:
				reg := valueToInt(inst.Tokens[5])
				op := 0xF01E | ((reg << 8) & 0x0F00)
				first := byte((op & 0xFF00) >> 8)
				second := byte(op & 0xFF)
				opcodes = append(opcodes, []byte{first, second}...)
			}
		case parse.SUB: // We're only going to implement VX = VX - VY
			reg1 := valueToInt(inst.Tokens[3])
			reg2 := valueToInt(inst.Tokens[8])
			op := 0x8005 | ((reg1 << 8) & 0x0F00) | ((reg2 << 4) & 0x00F0)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.OR:
			reg1 := valueToInt(inst.Tokens[3])
			reg2 := valueToInt(inst.Tokens[8])
			op := 0x8001 | ((reg1 << 8) & 0x0F00) | ((reg2 << 4) & 0x00F0)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.AND:
			reg1 := valueToInt(inst.Tokens[3])
			reg2 := valueToInt(inst.Tokens[8])
			op := 0x8002 | ((reg1 << 8) & 0x0F00) | ((reg2 << 4) & 0x00F0)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.XOR:
			reg1 := valueToInt(inst.Tokens[3])
			reg2 := valueToInt(inst.Tokens[8])
			op := 0x8003 | ((reg1 << 8) & 0x0F00) | ((reg2 << 4) & 0x00F0)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.SHR:
			reg := valueToInt(inst.Tokens[3])
			op := 0x8006 | ((reg << 8) & 0x0F00)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.SHL:
			reg := valueToInt(inst.Tokens[3])
			op := 0x800E | ((reg << 8) & 0x0F00)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.BRND:
			reg := valueToInt(inst.Tokens[3])
			val := valueToInt(inst.Tokens[6])
			op := 0xC000 | ((reg << 8) & 0x0F00) | (val & 0x00FF)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.DRW:
			reg1 := valueToInt(inst.Tokens[3])
			reg2 := valueToInt(inst.Tokens[8])
			op := 0xD000 | ((reg1 << 8) & 0x0F00) | ((reg2 << 4) & 0x00F0)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.FX29:
			reg := valueToInt(inst.Tokens[3])
			op := 0xF029 | ((reg << 8) & 0x0F00)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.FX33:
			reg := valueToInt(inst.Tokens[3])
			op := 0xF033 | ((reg << 8) & 0x0F00)
			first := byte((op & 0xFF00) >> 8)
			second := byte(op & 0xFF)
			opcodes = append(opcodes, []byte{first, second}...)
		case parse.FX55: // TODO: These last two need to be entirely worked throughout the system
			opcodes = append(opcodes, []byte{0xFF, 0x55}...)
		case parse.FX65:
			opcodes = append(opcodes, []byte{0xFF, 0x65}...)
		}
	}
	return opcodes
}

func valueToInt(token parse.Token) int {
	if token.Type == parse.HEX {
		val, err := strconv.ParseInt(token.Literal[2:], 16, 16)
		if err != nil {
			return -1
		}
		return int(val)
	} else {
		val, err := strconv.ParseInt(token.Literal, 10, 16)
		if err != nil {
			return -1
		}
		return int(val)
	}
}
