package compile

import (
	"chip8-emu/internal/chip8/parse"
	"fmt"
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
			bytes := ParseInstruction(0x00E0, inst)
			opcodes = append(opcodes, bytes...)
		case parse.SYSCALL:
			bytes := ParseInstruction(0x0000, inst)
			opcodes = append(opcodes, bytes...)
		case parse.CALL:
			bytes := ParseInstruction(0x2000, inst)
			opcodes = append(opcodes, bytes...)
		case parse.RET:
			bytes := ParseInstruction(0x00EE, inst)
			opcodes = append(opcodes, bytes...)
		case parse.JMP:
			bytes := ParseInstruction(0x1000, inst)
			opcodes = append(opcodes, bytes...)
		case parse.RJMP:
			bytes := ParseInstruction(0xB000, inst)
			opcodes = append(opcodes, bytes...)
		case parse.SEQ:
			switch inst.Format {
			case CMD_REG_VAL:
				bytes := ParseInstruction(0x3000, inst)
				opcodes = append(opcodes, bytes...)
			case CMD_REG_REG:
				bytes := ParseInstruction(0x5000, inst)
				opcodes = append(opcodes, bytes...)
			}
		case parse.SNEQ:
			switch inst.Format {
			case CMD_REG_VAL:
				bytes := ParseInstruction(0x4000, inst)
				opcodes = append(opcodes, bytes...)
			case CMD_REG_REG:
				bytes := ParseInstruction(0x9000, inst)
				opcodes = append(opcodes, bytes...)
			}
		case parse.JKP:
			bytes := ParseInstruction(0xE09E, inst)
			opcodes = append(opcodes, bytes...)
		case parse.JKNP:
			bytes := ParseInstruction(0xE0A1, inst)
			opcodes = append(opcodes, bytes...)
		case parse.WK:
			bytes := ParseInstruction(0xF00A, inst)
			opcodes = append(opcodes, bytes...)
		case parse.MOV:
			switch inst.Format {
			case CMD_REG_VAL:
				bytes := ParseInstruction(0x6000, inst)
				opcodes = append(opcodes, bytes...)
			case CMD_REG_REG:
				bytes := ParseInstruction(0x8000, inst)
				opcodes = append(opcodes, bytes...)
			case CMD_SPC_VAL:
				bytes := ParseInstruction(0xA000, inst)
				opcodes = append(opcodes, bytes...)
			case CMD_REG_SPC:
				bytes := ParseInstruction(0xF007, inst)
				opcodes = append(opcodes, bytes...)
			case CMD_SPC_REG:
				if inst.Tokens[1].Type == parse.DELAY {
					bytes := ParseInstruction(0xF015, inst)
					opcodes = append(opcodes, bytes...)
				} else { // SND_DELAY
					bytes := ParseInstruction(0xF018, inst)
					opcodes = append(opcodes, bytes...)
				}
			}
		case parse.ADD:
			switch inst.Format {
			case CMD_REG_VAL:
				bytes := ParseInstruction(0x7000, inst)
				opcodes = append(opcodes, bytes...)
			case CMD_REG_REG:
				bytes := ParseInstruction(0x8004, inst)
				opcodes = append(opcodes, bytes...)
			case CMD_SPC_REG:
				bytes := ParseInstruction(0xF01E, inst)
				opcodes = append(opcodes, bytes...)
			}
		case parse.SUB: // We're only going to implement VX = VX - VY
			bytes := ParseInstruction(0x8005, inst)
			opcodes = append(opcodes, bytes...)
		case parse.OR:
			bytes := ParseInstruction(0x8001, inst)
			opcodes = append(opcodes, bytes...)
		case parse.AND:
			bytes := ParseInstruction(0x8002, inst)
			opcodes = append(opcodes, bytes...)
		case parse.XOR:
			bytes := ParseInstruction(0x8003, inst)
			opcodes = append(opcodes, bytes...)
		case parse.SHR:
			bytes := ParseInstruction(0x8006, inst)
			opcodes = append(opcodes, bytes...)
		case parse.SHL:
			bytes := ParseInstruction(0x800E, inst)
			opcodes = append(opcodes, bytes...)
		case parse.BRND:
			bytes := ParseInstruction(0xC000, inst)
			opcodes = append(opcodes, bytes...)
		case parse.DRW:
			bytes := ParseInstruction(0xD000, inst)
			opcodes = append(opcodes, bytes...)
		case parse.FX29:
			bytes := ParseInstruction(0xF029, inst)
			opcodes = append(opcodes, bytes...)
		case parse.FX33:
			bytes := ParseInstruction(0xF033, inst)
			opcodes = append(opcodes, bytes...)
		case parse.FX55: // TODO: These last two need to be entirely worked throughout the system
			bytes := ParseInstruction(0xFF55, inst)
			opcodes = append(opcodes, bytes...)
		case parse.FX65:
			bytes := ParseInstruction(0xFF65, inst)
			opcodes = append(opcodes, bytes...)
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

func OpcodeToBytes(opcode int) []byte {
	first := byte(opcode >> 8)
	second := byte(opcode & 0xFF)
	return []byte{first, second}
}

func ParseInstruction(opcode int, instruction Instruction) []byte {
	switch instruction.Format {
	case SINGULAR:
		return ParseSINGULAR(opcode)
	case CMD_VAL:
		return OpcodeToBytes(ParseCMD_VAL(opcode, instruction))
	case CMD_REG:
		return OpcodeToBytes(ParseCMD_REG(opcode, instruction))
	case CMD_REG_VAL:
		return OpcodeToBytes(ParseCMD_REG_VAL(opcode, instruction))
	case CMD_REG_REG:
		return OpcodeToBytes(ParseCMD_REG_REG(opcode, instruction))
	case CMD_REG_SPC:
		return OpcodeToBytes(ParseCMD_REG_SPC(opcode, instruction))
	case CMD_SPC_REG:
		return OpcodeToBytes(ParseCMD_SPC_REG(opcode, instruction))
	case CMD_SPC_VAL:
		return OpcodeToBytes(ParseCMD_SPC_VAL(opcode, instruction))
	}
	fmt.Printf("This should never be reached! %d %#v\n", opcode, instruction)
	return []byte{}
}

func ParseSINGULAR(opcode int) []byte {
	return OpcodeToBytes(opcode)
}

func ParseCMD_VAL(opcode int, inst Instruction) int {
	val := valueToInt(inst.Tokens[1])
	op := opcode | (0xFFF & val)
	return op
}

func ParseCMD_REG(opcode int, inst Instruction) int {
	reg := valueToInt(inst.Tokens[3])
	op := opcode | ((reg << 8) & 0x0F00)
	return op
}

func ParseCMD_REG_VAL(opcode int, inst Instruction) int {
	reg := valueToInt(inst.Tokens[3])
	val := valueToInt(inst.Tokens[6])
	op := opcode | ((reg << 8) & 0x0F00) | (val & 0x00FF)
	return op
}

func ParseCMD_REG_REG(opcode int, inst Instruction) int {
	reg1 := valueToInt(inst.Tokens[3])
	reg2 := valueToInt(inst.Tokens[8])
	op := opcode | ((reg1 << 8) & 0x0F00) | ((reg2 << 4) & 0x00F0)
	return op
}

func ParseCMD_REG_SPC(opcode int, inst Instruction) int {
	reg := valueToInt(inst.Tokens[3])
	op := opcode | ((reg << 8) & 0x0F00)
	return op
}

func ParseCMD_SPC_REG(opcode int, inst Instruction) int {
	reg := valueToInt(inst.Tokens[5])
	op := opcode | ((reg << 8) & 0x0F00)
	return op
}

func ParseCMD_SPC_VAL(opcode int, inst Instruction) int {
	val := valueToInt(inst.Tokens[3])
	op := 0xA000 | (val & 0x0FFF)
	return op
}
