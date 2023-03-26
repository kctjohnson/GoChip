package compile

import (
	"fmt"

	"github.com/kctjohnson/chip8-emu/internal/chip8/parse"
)

type InstructionFormat int

const (
	SINGULAR InstructionFormat = iota
	CMD_VAL
	CMD_REG
	CMD_REG_VAL
	CMD_REG_REG
	CMD_REG_SPC
	CMD_SPC_REG
	CMD_SPC_VAL
)

type Instruction struct {
	Format InstructionFormat
	Tokens []parse.Token
}

type InstructionSet struct {
	Instructions []Instruction
}

func NewInstructionSet(tokens []parse.Token) *InstructionSet {
	is := &InstructionSet{}
	is.parse(tokens)
	return is
}

func (is *InstructionSet) parse(tokens []parse.Token) {
	is.Instructions = []Instruction{}
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
			is.Instructions = append(is.Instructions, Instruction{
				Format: SINGULAR,
				Tokens: []parse.Token{curToken},
			})
		case parse.SYSCALL:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_VAL,
				Tokens: tokens[i : i+2],
			})
			i++
		case parse.CALL:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_VAL,
				Tokens: tokens[i : i+2],
			})
			i++
		case parse.RET:
			is.Instructions = append(is.Instructions, Instruction{
				Format: SINGULAR,
				Tokens: []parse.Token{curToken},
			})
		case parse.JMP:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_VAL,
				Tokens: tokens[i : i+2],
			})
			i++
		case parse.RJMP:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_VAL,
				Tokens: tokens[i : i+2],
			})
			i++
		case parse.SEQ:
			if tokens[i+6].Type == parse.REG {
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_REG_REG,
					Tokens: tokens[i : i+10],
				})
				i += 9
			} else {
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_REG_VAL,
					Tokens: tokens[i : i+7],
				})
				i += 6
			}
		case parse.SNEQ:
			if tokens[i+6].Type == parse.REG {
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_REG_REG,
					Tokens: tokens[i : i+10],
				})
				i += 9
			} else {
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_REG_VAL,
					Tokens: tokens[i : i+7],
				})
				i += 6
			}
		case parse.JKP:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parse.JKNP:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parse.WK:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parse.MOV:
			switch tokens[i+1].Type {
			case parse.REG:
				if tokens[i+6].Type == parse.REG {
					is.Instructions = append(is.Instructions, Instruction{
						Format: CMD_REG_REG,
						Tokens: tokens[i : i+10],
					})
					i += 9
				} else if tokens[i+6].Type == parse.DELAY {
					is.Instructions = append(is.Instructions, Instruction{
						Format: CMD_REG_SPC,
						Tokens: tokens[i : i+7],
					})
					i += 6
				} else {
					is.Instructions = append(is.Instructions, Instruction{
						Format: CMD_REG_VAL,
						Tokens: tokens[i : i+7],
					})
					i += 6
				}
			case parse.ADP:
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_SPC_VAL,
					Tokens: tokens[i : i+4],
				})
				i += 3
			case parse.DELAY:
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_SPC_REG,
					Tokens: tokens[i : i+7],
				})
				i += 6
			case parse.SND_DELAY:
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_SPC_REG,
					Tokens: tokens[i : i+7],
				})
				i += 6
			}
		case parse.ADD:
			switch tokens[i+1].Type {
			case parse.REG:
				switch tokens[i+6].Type {
				case parse.REG:
					is.Instructions = append(is.Instructions, Instruction{
						Format: CMD_REG_REG,
						Tokens: tokens[i : i+10],
					})
					i += 9
				default:
					is.Instructions = append(is.Instructions, Instruction{
						Format: CMD_REG_VAL,
						Tokens: tokens[i : i+7],
					})
					i += 6
				}
			case parse.ADP:
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_SPC_REG,
					Tokens: tokens[i : i+7],
				})
				i += 6
			}
		case parse.SUB:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG_REG,
				Tokens: tokens[i : i+10],
			})
			i += 9
		case parse.OR:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG_REG,
				Tokens: tokens[i : i+10],
			})
			i += 9
		case parse.AND:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG_REG,
				Tokens: tokens[i : i+10],
			})
			i += 9
		case parse.XOR:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG_REG,
				Tokens: tokens[i : i+10],
			})
			i += 9
		case parse.SHR:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parse.SHL:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parse.BRND:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+7],
			})
			i += 4
		case parse.DRW:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG_REG,
				Tokens: tokens[i : i+10],
			})
			i += 9
		case parse.FX29:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parse.FX33:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parse.FX55:
			is.Instructions = append(is.Instructions, Instruction{
				Format: SINGULAR,
				Tokens: []parse.Token{curToken},
			})
		case parse.FX65:
			is.Instructions = append(is.Instructions, Instruction{
				Format: SINGULAR,
				Tokens: []parse.Token{curToken},
			})
		}
	}
}
