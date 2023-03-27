package disassembler

import (
	"fmt"

	"github.com/kctjohnson/chip8-emu/internal/chip8"
	"github.com/kctjohnson/chip8-emu/internal/chip8/emulator"
)

func Disassemble(gameFilePath string) {
	emu := emulator.NewEmulator(gameFilePath)

	zeroCounter := 0
	for zeroCounter < 5 {
		op := emu.GetNextOpcode()

		fmt.Print(DisassembleOpcode(op))

		if op == 0 {
			zeroCounter++
		} else {
			zeroCounter = 0
		}
	}
}

func DisassembleOpcode(op chip8.WORD) string {
	str := fmt.Sprintf("0x%04X  ", op)
	switch op & 0xF000 {
	case 0x0000:
		switch op & 0x00FF {
		case 0x00E0:
			str += "CLS\n"
		case 0x00EE:
			str += "RET\n"
		default:
			str += fmt.Sprintf("SYSCALL 0x%04X\n", op&0x0FFF)
		}
	case 0x1000:
		str += fmt.Sprintf("JMP 0x%04X\n", op&0x0FFF)
	case 0x2000:
		str += fmt.Sprintf("CALL 0x%04X\n", op&0x0FFF)
	case 0x3000:
		regx, _ := chip8.GetXYReg(op)
		nn := op & 0x00FF
		str += fmt.Sprintf("SEQ reg[0x%X], %d\n", regx, nn)
	case 0x4000:
		regx, _ := chip8.GetXYReg(op)
		nn := op & 0x00FF
		str += fmt.Sprintf("SNEQ reg[0x%X], %d\n", regx, nn)
	case 0x5000:
		regx, regy := chip8.GetXYReg(op)
		str += fmt.Sprintf("SEQ reg[0x%X], reg[0x%X]\n", regx, regy)
	case 0x6000:
		regx, _ := chip8.GetXYReg(op)
		nn := op & 0x00FF
		str += fmt.Sprintf("MOV reg[0x%X], %d\n", regx, nn)
	case 0x7000:
		regx, _ := chip8.GetXYReg(op)
		nn := op & 0x00FF
		str += fmt.Sprintf("ADD reg[0x%X], %d\n", regx, nn)
	case 0x8000:
		switch op & 0x000F {
		case 0x0:
			regx, regy := chip8.GetXYReg(op)
			str += fmt.Sprintf("MOV reg[0x%X], reg[0x%X]\n", regx, regy)
		case 0x1:
			regx, regy := chip8.GetXYReg(op)
			str += fmt.Sprintf("OR reg[0x%X], reg[0x%X]\n", regx, regy)
		case 0x2:
			regx, regy := chip8.GetXYReg(op)
			str += fmt.Sprintf("AND reg[0x%X], reg[0x%X]\n", regx, regy)
		case 0x3:
			regx, regy := chip8.GetXYReg(op)
			str += fmt.Sprintf("XOR reg[0x%X], reg[0x%X]\n", regx, regy)
		case 0x4:
			regx, regy := chip8.GetXYReg(op)
			str += fmt.Sprintf("ADD reg[0x%X], reg[0x%X]\n", regx, regy)
		case 0x5:
			regx, regy := chip8.GetXYReg(op)
			str += fmt.Sprintf("SUB reg[0x%X], reg[0x%X]\n", regx, regy)
		case 0x6:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("SHR reg[0x%X]\n", regx)
		case 0x7:
			regx, regy := chip8.GetXYReg(op)
			str += fmt.Sprintf("SUB reg[0x%X], reg[0x%X]\n", regx, regy)
		case 0xE:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("SHL reg[0x%X]\n", regx)
		default:
			str += "Something went wrong!\n"
		}
	case 0x9000:
		regx, regy := chip8.GetXYReg(op)
		str += fmt.Sprintf("SNEQ reg[0x%X], reg[0x%X]\n", regx, regy)
	case 0xA000:
		str += fmt.Sprintf("MOV I, %d\n", op&0x0FFF)
	case 0xB000:
		str += fmt.Sprintf("RJMP 0x%X\n", op&0x0FFF)
	case 0xC000:
		regx, _ := chip8.GetXYReg(op)
		str += fmt.Sprintf("BRND reg[0x%X], %d\n", regx, op&0x00FF)
	case 0xD000:
		regx, regy := chip8.GetXYReg(op)
		rows := op & 0xF
		str += fmt.Sprintf("DRW reg[0x%X], reg[0x%X], %d\n", regx, regy, rows)
	case 0xE000:
		switch op & 0x00FF {
		case 0x9E:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("JKP reg[0x%X]\n", regx)
		case 0xA1:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("JKNP reg[0x%X]\n", regx)
		default:
			str += "Something went wrong!\n"
		}
	case 0xF000:
		switch op & 0x00FF {
		case 0x07:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("MOV reg[0x%X], DELAY\n", regx)
		case 0x0A:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("WK reg[0x%X]\n", regx)
		case 0x15:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("MOV DELAY, reg[0x%X]\n", regx)
		case 0x18:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("MOV SND_DELAY, reg[0x%X]\n", regx)
		case 0x1E:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("ADD I, reg[0x%X]\n", regx)
		case 0x29:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("FX29 reg[0x%X]\n", regx)
		case 0x33:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("FX33 reg[0x%X]\n", regx)
		case 0x55:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("FX55 reg[0x%X]\n", regx)
		case 0x65:
			regx, _ := chip8.GetXYReg(op)
			str += fmt.Sprintf("FX65 reg[0x%X]\n", regx)
		default:
			str += "Something went wrong!\n"
		}
	default:
		str += "Something went wrong!\n"
	}
	return str
}
