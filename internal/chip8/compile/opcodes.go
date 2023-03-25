package compile

import (
	"chip8-emu/internal/chip8/parse"
	"fmt"
)

type InstructionFormat string

const (
	SINGULAR    InstructionFormat = "CMD"
	CMD_VAL     InstructionFormat = "CMD VAL"
	CMD_REG     InstructionFormat = "CMD REG"
	CMD_ADP     InstructionFormat = "CMD ADP"
	CMD_REG_VAL InstructionFormat = "CMD REG, VAL"
	CMD_REG_REG InstructionFormat = "CMD REG, REG"
	CMD_REG_DLY InstructionFormat = "CMD REG, DLY"
	CMD_DLY_REG InstructionFormat = "CMD DLY, REG"
	CMD_ADP_REG InstructionFormat = "CMD ADP, REG"
)

type Instruction struct {
	format InstructionFormat
	tokens []parse.Token
}

type InstructionSet struct {
	instructions []Instruction
}

func NewInstructionSet(tokens []parse.Token) *InstructionSet {
	is := &InstructionSet{}
	is.parse(tokens)
	return is
}

func (is *InstructionSet) parse(tokens []parse.Token) {
	is.instructions = []Instruction{}
	for i := 0; i < len(tokens); i++ {
		curToken := tokens[i]
		switch curToken.Type {
		case parse.ILLEGAL:
			fmt.Printf("Illegal token found! %#v\n", curToken)
			return
		case parse.EOF:
			return
		case parse.UNKNOWNIDENT:
			fmt.Printf("Unknown identifier found! %#v\n", curToken)
			return
		case parse.CLS:
			is.instructions = append(is.instructions, Instruction{
				format: SINGULAR,
				tokens: []parse.Token{curToken},
			})
		case parse.SYSCALL:
			is.instructions = append(is.instructions, Instruction{
				format: CMD_VAL,
				tokens: tokens[i : i+2],
			})
			i++
		case parse.CALL:
			is.instructions = append(is.instructions, Instruction{
				format: CMD_VAL,
				tokens: tokens[i : i+2],
			})
			i++
		case parse.RET:
			is.instructions = append(is.instructions, Instruction{
				format: SINGULAR,
				tokens: []parse.Token{curToken},
			})
		case parse.JMP:
			is.instructions = append(is.instructions, Instruction{
				format: CMD_VAL,
				tokens: tokens[i : i+2],
			})
			i++
		case parse.RJMP:
			is.instructions = append(is.instructions, Instruction{
				format: CMD_VAL,
				tokens: tokens[i : i+2],
			})
			i++
		case parse.SEQ:
			if tokens[i+6].Type == parse.REG {
				is.instructions = append(is.instructions, Instruction{
					format: CMD_REG_REG,
					tokens: tokens[i:10],
				})
				i += 9
			} else {
				is.instructions = append(is.instructions, Instruction{
					format: CMD_REG_VAL,
					tokens: tokens[i:7],
				})
				i += 6
			}
		case parse.SNEQ:
		case parse.JKP:
		case parse.JKNP:
		case parse.WK:
		case parse.MOV:
		case parse.ADD:
		case parse.SUB:
		case parse.OR:
		case parse.AND:
		case parse.XOR:
		case parse.SHR:
		case parse.SHL:
		case parse.BRND:
		case parse.DRW:
		case parse.FX29:
		case parse.FX33:
		case parse.FX55:
		case parse.FX65:
		}
	}
}
