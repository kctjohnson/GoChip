package compiler

import (
	"fmt"

	"github.com/kctjohnson/chip8-emu/internal/chip8/parser"
)

type InstructionFormat int

const (
	CMD InstructionFormat = iota
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
	Tokens []parser.Token
}

type InstructionSet struct {
	Instructions []Instruction
}

func NewInstructionSet(tokens []parser.Token) *InstructionSet {
	is := &InstructionSet{}
	is.parse(tokens)
	return is
}

func (is *InstructionSet) parse(tokens []parser.Token) {
	is.Instructions = []Instruction{}
	for i := 0; i < len(tokens); i++ {
		curToken := tokens[i]
		switch curToken.Type {
		case parser.ILLEGAL:
			fmt.Printf("Illegal token found! %#v\n", curToken)
			return
		case parser.EOF:
			return
		case parser.UNKNOWNIDENT:
			fmt.Printf("Unknown identifier found! %#v\n", curToken)
			return
		case parser.CLS:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD,
				Tokens: []parser.Token{curToken},
			})
		case parser.SYSCALL:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_VAL,
				Tokens: tokens[i : i+2],
			})
			i++
		case parser.CALL:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_VAL,
				Tokens: tokens[i : i+2],
			})
			i++
		case parser.RET:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD,
				Tokens: []parser.Token{curToken},
			})
		case parser.JMP:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_VAL,
				Tokens: tokens[i : i+2],
			})
			i++
		case parser.RJMP:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_VAL,
				Tokens: tokens[i : i+2],
			})
			i++
		case parser.SEQ:
			if tokens[i+6].Type == parser.REG {
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
		case parser.SNEQ:
			if tokens[i+6].Type == parser.REG {
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
		case parser.JKP:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parser.JKNP:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parser.WK:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parser.MOV:
			switch tokens[i+1].Type {
			case parser.REG:
				if tokens[i+6].Type == parser.REG {
					is.Instructions = append(is.Instructions, Instruction{
						Format: CMD_REG_REG,
						Tokens: tokens[i : i+10],
					})
					i += 9
				} else if tokens[i+6].Type == parser.DELAY {
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
			case parser.ADP:
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_SPC_VAL,
					Tokens: tokens[i : i+4],
				})
				i += 3
			case parser.DELAY:
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_SPC_REG,
					Tokens: tokens[i : i+7],
				})
				i += 6
			case parser.SND_DELAY:
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_SPC_REG,
					Tokens: tokens[i : i+7],
				})
				i += 6
			}
		case parser.ADD:
			switch tokens[i+1].Type {
			case parser.REG:
				switch tokens[i+6].Type {
				case parser.REG:
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
			case parser.ADP:
				is.Instructions = append(is.Instructions, Instruction{
					Format: CMD_SPC_REG,
					Tokens: tokens[i : i+7],
				})
				i += 6
			}
		case parser.SUB:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG_REG,
				Tokens: tokens[i : i+10],
			})
			i += 9
		case parser.OR:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG_REG,
				Tokens: tokens[i : i+10],
			})
			i += 9
		case parser.AND:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG_REG,
				Tokens: tokens[i : i+10],
			})
			i += 9
		case parser.XOR:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG_REG,
				Tokens: tokens[i : i+10],
			})
			i += 9
		case parser.SHR:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parser.SHL:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parser.BRND:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+7],
			})
			i += 4
		case parser.DRW:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG_REG,
				Tokens: tokens[i : i+10],
			})
			i += 9
		case parser.FX29:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parser.FX33:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
			i += 4
		case parser.FX55:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
		case parser.FX65:
			is.Instructions = append(is.Instructions, Instruction{
				Format: CMD_REG,
				Tokens: tokens[i : i+5],
			})
		}
	}
}
