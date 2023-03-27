package compiler

import (
	"fmt"
	"strconv"

	"github.com/kctjohnson/chip8-emu/internal/chip8/parser"
)

type Compiler struct {
	Instructions *InstructionSet
	tokens       []parser.Token
}

func NewCompiler(tokens []parser.Token) *Compiler {
	return &Compiler{
		tokens:       tokens,
		Instructions: NewInstructionSet(tokens),
	}
}

type mapKey struct {
	Type   parser.TokenType
	Format InstructionFormat
}

var OpcodeMap = map[mapKey]int{
	{Type: parser.CLS, Format: CMD}:             0x00E0,
	{Type: parser.SYSCALL, Format: CMD_VAL}:     0x0000,
	{Type: parser.CALL, Format: CMD_VAL}:        0x2000,
	{Type: parser.RET, Format: CMD}:             0x00EE,
	{Type: parser.JMP, Format: CMD_VAL}:         0x1000,
	{Type: parser.RJMP, Format: CMD_VAL}:        0xB000,
	{Type: parser.SEQ, Format: CMD_REG_REG}:     0x5000,
	{Type: parser.SEQ, Format: CMD_REG_VAL}:     0x3000,
	{Type: parser.SNEQ, Format: CMD_REG_REG}:    0x9000,
	{Type: parser.SNEQ, Format: CMD_REG_VAL}:    0x4000,
	{Type: parser.JKP, Format: CMD_REG}:         0xE09E,
	{Type: parser.JKNP, Format: CMD_REG}:        0xE0A1,
	{Type: parser.WK, Format: CMD_REG}:          0xF00A,
	{Type: parser.MOV, Format: CMD_REG_REG}:     0x8000,
	{Type: parser.MOV, Format: CMD_REG_SPC}:     0xF007,
	{Type: parser.MOV, Format: CMD_REG_VAL}:     0x6000,
	{Type: parser.MOV, Format: CMD_SPC_VAL}:     0xA000,
	{Type: parser.ADD, Format: CMD_REG_REG}:     0x8004,
	{Type: parser.ADD, Format: CMD_REG_VAL}:     0x7000,
	{Type: parser.ADD, Format: CMD_SPC_REG}:     0xF01E,
	{Type: parser.SUB, Format: CMD_REG_REG}:     0x8005,
	{Type: parser.OR, Format: CMD_REG_REG}:      0x8001,
	{Type: parser.AND, Format: CMD_REG_REG}:     0x8002,
	{Type: parser.XOR, Format: CMD_REG_REG}:     0x8003,
	{Type: parser.SHR, Format: CMD_REG}:         0x8006,
	{Type: parser.SHL, Format: CMD_REG}:         0x800E,
	{Type: parser.BRND, Format: CMD_REG}:        0xC000,
	{Type: parser.DRW, Format: CMD_REG_REG_VAL}: 0xD000,
	{Type: parser.FX29, Format: CMD_REG}:        0xF029,
	{Type: parser.FX33, Format: CMD_REG}:        0xF033,
	{Type: parser.FX55, Format: CMD_REG}:        0xF055,
	{Type: parser.FX65, Format: CMD_REG}:        0xF065,
}

func (c Compiler) Compile() []byte {
	opcodes := []byte{}
	for _, inst := range c.Instructions.Instructions {
		// MOV is a special case due to sound delay and time delay
		if inst.Tokens[0].Type == parser.MOV && inst.Format == CMD_SPC_REG {
			if inst.Tokens[1].Type == parser.DELAY {
				bytes := ParseInstruction(0xF015, inst)
				opcodes = append(opcodes, bytes...)
			} else { // SND_DELAY
				bytes := ParseInstruction(0xF018, inst)
				opcodes = append(opcodes, bytes...)
			}
		} else {
			opValue, ok := OpcodeMap[mapKey{
				Type:   inst.Tokens[0].Type,
				Format: inst.Format,
			}]

			if !ok {
				panic(fmt.Sprintf("Invalid instruction! %#v\n", inst))
			}

			bytes := ParseInstruction(opValue, inst)
			opcodes = append(opcodes, bytes...)
		}
	}
	return opcodes
}

func valueToInt(token parser.Token) int {
	if token.Type == parser.HEX {
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
	case CMD:
		return ParseCMD(opcode)
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
	case CMD_REG_REG_VAL:
		return OpcodeToBytes(ParseCMD_REG_REG_VAL(opcode, instruction))
	}
	fmt.Printf("This should never be reached! %d %#v\n", opcode, instruction)
	return []byte{}
}

func ParseCMD(opcode int) []byte {
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

func ParseCMD_REG_REG_VAL(opcode int, inst Instruction) int {
	// DRW REG[0x1], REG[0x2], N
	reg1 := valueToInt(inst.Tokens[3])
	reg2 := valueToInt(inst.Tokens[8])
	val := valueToInt(inst.Tokens[11])
	op := opcode | ((reg1 << 8) & 0x0F00) | ((reg2 << 4) & 0x00F0) | (val & 0xF)
	return op
}
