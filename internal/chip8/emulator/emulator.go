package emulator

import (
	"chip8-emu/internal/chip8"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	debug               = false
	timerCap chip8.BYTE = 60
)

type Emulator struct {
	Memory        [0xFFF]chip8.BYTE
	Registers     [16]chip8.BYTE
	I             chip8.WORD
	PC            chip8.WORD
	Stack         []chip8.WORD
	Inputs        [16]chip8.BYTE
	Delay         chip8.BYTE
	SoundDelay    chip8.BYTE
	CurrentOpcode chip8.WORD

	FilePath string

	LastTick     time.Time
	DelayEnabled bool

	ScreenData [64][32]chip8.BYTE
}

func NewEmulator(gameFilePath string) *Emulator {
	emu := &Emulator{
		FilePath: gameFilePath,
	}
	emu.CPUReset()
	return emu
}

func (h *Emulator) CPUReset() {
	// Reset all of the data
	h.I = 0
	h.PC = 0x200
	for i := range h.Registers {
		h.Registers[i] = 0
	}
	h.DelayEnabled = true
	h.Delay = 60
	h.LastTick = time.Now()
	h.Stack = []chip8.WORD{}
	h.Inputs = [16]chip8.BYTE{}
	h.ScreenData = [64][32]chip8.BYTE{}
	h.Memory = [0xFFF]chip8.BYTE{}

	// Load in the game
	gameFile, err := os.ReadFile(h.FilePath)
	if err != nil {
		panic(err)
	}

	for i := range gameFile {
		h.Memory[i+0x200] = chip8.BYTE(gameFile[i])
	}
}

func (h *Emulator) GetNextOpcode() chip8.WORD {
	res := (chip8.WORD(h.Memory[h.PC]) << 8) | chip8.WORD(h.Memory[h.PC+1])
	h.CurrentOpcode = res
	h.PC += 2
	return res
}

func (h Emulator) GetOpcode(pc chip8.WORD) chip8.WORD {
	return (chip8.WORD(h.Memory[pc]) << 8) | chip8.WORD(h.Memory[pc+1])
}

func (h *Emulator) Step() {
	op := h.GetNextOpcode()

	switch op & 0xF000 {
	case 0x0000:
		switch op & 0x00FF {
		case 0x00E0:
			h.Opcode00E0(op)
		case 0x00EE:
			h.Opcode00EE(op)
		default:
			h.Opcode0NNN(op)
		}
	case 0x1000:
		h.Opcode1NNN(op)
	case 0x2000:
		h.Opcode2NNN(op)
	case 0x3000:
		h.Opcode3XNN(op)
	case 0x4000:
		h.Opcode4XNN(op)
	case 0x5000:
		h.Opcode5XY0(op)
	case 0x6000:
		h.Opcode6XNN(op)
	case 0x7000:
		h.Opcode7XNN(op)
	case 0x8000:
		switch op & 0x000F {
		case 0x0:
			h.Opcode8XY0(op)
		case 0x1:
			h.Opcode8XY1(op)
		case 0x2:
			h.Opcode8XY2(op)
		case 0x3:
			h.Opcode8XY3(op)
		case 0x4:
			h.Opcode8XY4(op)
		case 0x5:
			h.Opcode8XY5(op)
		case 0x6:
			h.Opcode8XY6(op)
		case 0x7:
			h.Opcode8XY7(op)
		case 0xE:
			h.Opcode8XYE(op)
		}
	case 0x9000:
		h.Opcode9XY0(op)
	case 0xA000:
		h.OpcodeANNN(op)
	case 0xB000:
		h.OpcodeBNNN(op)
	case 0xC000:
		h.OpcodeCXNN(op)
	case 0xD000:
		h.OpcodeDXYN(op)
	case 0xE000:
		switch op & 0x00FF {
		case 0x9E:
			h.OpcodeEX9E(op)
		case 0xA1:
			h.OpcodeEXA1(op)
		}
	case 0xF000:
		switch op & 0x00FF {
		case 0x07:
			h.OpcodeFX07(op)
		case 0x0A:
			h.OpcodeFX0A(op)
		case 0x15:
			h.OpcodeFX15(op)
		case 0x1E:
			h.OpcodeFX1E(op)
		case 0x29:
			h.OpcodeFX29(op)
		case 0x33:
			h.OpcodeFX33(op)
		case 0x55:
			h.OpcodeFX55(op)
		case 0x65:
			h.OpcodeFX65(op)
		}
	default:
		fmt.Printf("Something went wrong!\n")
	}

	if h.DelayEnabled {
		elapsedTime := time.Since(h.LastTick)
		tickSpeed := time.Second / time.Duration(timerCap)
		if elapsedTime > tickSpeed {
			delta := elapsedTime - tickSpeed
			h.LastTick = time.Now().Add(-delta)
			h.Delay -= 1
			if h.Delay > timerCap {
				h.Delay = timerCap
			}
		}
	} else {
		if h.Delay > 0 {
			elapsedTime := time.Since(h.LastTick)
			tickSpeed := time.Second / time.Duration(timerCap)
			if elapsedTime > tickSpeed {
				delta := elapsedTime - tickSpeed
				h.LastTick = time.Now().Add(-delta)
				h.Delay -= 1
			}
		}
	}
}

// Call machine code routine at NNN
func (h *Emulator) Opcode0NNN(op chip8.WORD) {
	if debug {
		fmt.Printf("0NNN: %x\n", op)
	}
	fmt.Println("0NNN Not implemented!")
}

// Clear the screen
func (h *Emulator) Opcode00E0(op chip8.WORD) {
	if debug {
		fmt.Printf("00E0: %x\n", op)
	}
	for x := range h.ScreenData {
		for y := range h.ScreenData[x] {
			h.ScreenData[x][y] = 0
		}
	}
}

// Return from a subroutine
func (h *Emulator) Opcode00EE(op chip8.WORD) {
	if debug {
		fmt.Printf("00EE: 00%x\n", op)
	}
	returnAddress := h.Stack[len(h.Stack)-1]
	h.Stack = h.Stack[:len(h.Stack)-1]
	h.PC = returnAddress
}

// Jump to address NNN
func (h *Emulator) Opcode1NNN(op chip8.WORD) {
	if debug {
		fmt.Printf("1NNN: %x\n", op)
	}
	h.PC = op & 0x0FFF
}

// Call subroutine at NNN
func (h *Emulator) Opcode2NNN(op chip8.WORD) {
	if debug {
		fmt.Printf("2NNN: %x\n", op)
	}
	h.Stack = append(h.Stack, h.PC) // Save the program counter
	h.PC = op & 0x0FFF              // Jump to address NNN
}

// Skips the next instruction if VX equals NN
func (h *Emulator) Opcode3XNN(op chip8.WORD) {
	if debug {
		fmt.Printf("3XNN: %x\n", op)
	}
	regx := (op & 0x0F00) >> 8
	nn := op & 0x00FF
	if chip8.WORD(h.Registers[regx]) == nn {
		h.PC += 2
	}
}

// Skips the next instruction if VX does not equal NN
func (h *Emulator) Opcode4XNN(op chip8.WORD) {
	if debug {
		fmt.Printf("4XNN: %x\n", op)
	}
	regx := (op & 0x0F00) >> 8
	nn := op & 0x00FF
	if chip8.WORD(h.Registers[regx]) != nn {
		h.PC += 2
	}
}

// Skips the next instruction if VX equals VY
func (h *Emulator) Opcode5XY0(op chip8.WORD) {
	if debug {
		fmt.Printf("5XY0: %x\n", op)
	}
	regx, regy := chip8.GetXYReg(op)
	if h.Registers[regx] == h.Registers[regy] {
		h.PC += 2
	}
}

// Sets VX to NN
func (h *Emulator) Opcode6XNN(op chip8.WORD) {
	if debug {
		fmt.Printf("6XNN: %x\n", op)
	}
	regx := (op & 0x0F00) >> 8
	h.Registers[regx] = chip8.BYTE(op & 0x00FF)
}

// Adds NN to VX (carry flag is not changed)
func (h *Emulator) Opcode7XNN(op chip8.WORD) {
	if debug {
		fmt.Printf("7XNN: %x\n", op)
	}
	regx := (op & 0x0F00) >> 8
	h.Registers[regx] += chip8.BYTE(op & 0x00FF)
}

// Sets VX to the value of VY
func (h *Emulator) Opcode8XY0(op chip8.WORD) {
	if debug {
		fmt.Printf("8XY0: %x\n", op)
	}
	regx, regy := chip8.GetXYReg(op)
	h.Registers[regx] = h.Registers[regy]
}

// Sets VX to VX or VY (bitwise OR operation)
func (h *Emulator) Opcode8XY1(op chip8.WORD) {
	if debug {
		fmt.Printf("8XY1: %x\n", op)
	}
	regx, regy := chip8.GetXYReg(op)
	h.Registers[regx] = h.Registers[regx] | h.Registers[regy]
}

// Sets VX to VX and VY (bitwise AND operation)
func (h *Emulator) Opcode8XY2(op chip8.WORD) {
	if debug {
		fmt.Printf("8XY2: %x\n", op)
	}
	regx, regy := chip8.GetXYReg(op)
	h.Registers[regx] = h.Registers[regx] & h.Registers[regy]
}

// Sets VX to VX xor VY
func (h *Emulator) Opcode8XY3(op chip8.WORD) {
	if debug {
		fmt.Printf("8XY3: %x\n", op)
	}
	regx, regy := chip8.GetXYReg(op)
	h.Registers[regx] = h.Registers[regx] ^ h.Registers[regy]
}

// Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there is not
func (h *Emulator) Opcode8XY4(op chip8.WORD) {
	if debug {
		fmt.Printf("8XY4: %x\n", op)
	}
	regx, regy := chip8.GetXYReg(op)
	if h.Registers[regx]+h.Registers[regy] > 255 {
		h.Registers[0xF] = 1
	} else {
		h.Registers[0xF] = 0
	}
	h.Registers[regx] += h.Registers[regy]
}

// VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there is not
func (h *Emulator) Opcode8XY5(op chip8.WORD) {
	if debug {
		fmt.Printf("8XY5: %x\n", op)
	}
	h.Registers[0xF] = 1
	regx, regy := chip8.GetXYReg(op)
	xVal := h.Registers[regx]
	yVal := h.Registers[regy]
	if yVal > xVal { // If this is true will result in a value < 0
		h.Registers[0xF] = 0
	}
	h.Registers[regx] = xVal - yVal
}

// Stores the least significant bit of VX in VF and then shifts VX to the right by 1
func (h *Emulator) Opcode8XY6(op chip8.WORD) {
	if debug {
		fmt.Printf("8XY6: %x\n", op)
	}
	regx, _ := chip8.GetXYReg(op)
	h.Registers[0xF] = h.Registers[regx] & 1
	h.Registers[regx] >>= 1
}

// Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there is not
func (h *Emulator) Opcode8XY7(op chip8.WORD) {
	if debug {
		fmt.Printf("8XY7: %x\n", op)
	}
	regx, regy := chip8.GetXYReg(op)
	if h.Registers[regx] > h.Registers[regy] {
		h.Registers[0xF] = 0
	} else {
		h.Registers[0xF] = 1
	}
	h.Registers[regx] = h.Registers[regy] - h.Registers[regx]
}

// Stores the most significant bit of VX in VF and then shifts VX to the left by 1
func (h *Emulator) Opcode8XYE(op chip8.WORD) {
	if debug {
		fmt.Printf("8XYE: %x\n", op)
	}
	regx, _ := chip8.GetXYReg(op)
	h.Registers[0xF] = h.Registers[regx] & 0x80
	h.Registers[regx] <<= 1
}

// Skips the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block);
func (h *Emulator) Opcode9XY0(op chip8.WORD) {
	if debug {
		fmt.Printf("9XY0: %x\n", op)
	}
	regx, regy := chip8.GetXYReg(op)
	if h.Registers[regx] != h.Registers[regy] {
		h.PC += 2
	}
}

// Sets I to the address NNN.
func (h *Emulator) OpcodeANNN(op chip8.WORD) {
	if debug {
		fmt.Printf("ANNN: %x\n", op)
	}
	h.I = op & 0x0FFF
}

// Jumps to the address NNN plus V0.
func (h *Emulator) OpcodeBNNN(op chip8.WORD) {
	if debug {
		fmt.Printf("BNNN: %x\n", op)
	}
	h.PC = (op & 0x0FFF) + chip8.WORD(h.Registers[0])
}

// Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
func (h *Emulator) OpcodeCXNN(op chip8.WORD) {
	if debug {
		fmt.Printf("CXNN: %x\n", op)
	}
	h.Registers[(op&0x0F00)>>8] = chip8.BYTE(int((op & 0x00FF)) & rand.Intn(254))
}

// Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels. Each row of 8 pixels is read as bit-coded starting from memory location I; I value does not change after the execution of this instruction. As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that does not happen.
func (h *Emulator) OpcodeDXYN(op chip8.WORD) {
	if debug {
		fmt.Printf("DXYN: %x\n", op)
	}
	regx, regy := chip8.GetXYReg(op)

	height := op & 0x000F
	coordx := h.Registers[regx]
	coordy := h.Registers[regy]

	h.Registers[0xF] = 0

	// Clear out any leftovers
	for x := range h.ScreenData {
		for y := range h.ScreenData[x] {
			if h.ScreenData[x][y] == 2 {
				h.ScreenData[x][y] = 0
			}
		}
	}

	// Loop for the amount of vertical lines needed to draw this
	var yline chip8.WORD
	for yline = 0; yline < height && yline < 32; yline++ {
		data := h.Memory[h.I+yline]
		xpixelinv := 7
		xpixel := 0
		for xpixel = 0; xpixel < 8 && xpixel+int(coordx) < 64; {
			mask := chip8.BYTE(1 << xpixelinv)
			if data&mask > 0 {
				x := int(coordx) + xpixel
				y := chip8.WORD(coordy) + yline

				if x < 64 && y < 32 {
					if h.ScreenData[x][y] == 1 {
						h.Registers[0xF] = 1 // Collision
						h.ScreenData[x][y] = 2
					} else {
						h.ScreenData[x][y] ^= 1
					}
				}
			}
			xpixel++
			xpixelinv--
		}
	}
}

// Skips the next instruction if the key stored in VX is pressed (usually the next instruction is a jump to skip a code block).
func (h *Emulator) OpcodeEX9E(op chip8.WORD) {
	if debug {
		fmt.Printf("EX9E: %x\n", op)
	}
	regx := (op & 0x0F00) >> 8
	if h.Inputs[h.Registers[regx]] != 0 {
		h.PC += 2
		h.Inputs[h.Registers[regx]] = 0
	}
}

// Skips the next instruction if the key stored in VX is not pressed (usually the next instruction is a jump to skip a code block).
func (h *Emulator) OpcodeEXA1(op chip8.WORD) {
	if debug {
		fmt.Printf("EXA1: %x\n", op)
	}
	regx := (op & 0x0F00) >> 8
	if h.Inputs[h.Registers[regx]] == 0 {
		h.PC += 2
	} else {
		h.Inputs[h.Registers[regx]] = 0
	}
}

// Sets VX to the value of the delay timer.
func (h *Emulator) OpcodeFX07(op chip8.WORD) {
	if debug {
		fmt.Printf("FX07: %x\n", op)
	}
	regx, _ := chip8.GetXYReg(op)
	h.Registers[regx] = h.Delay
}

// A key press is awaited, and then stored in VX (blocking operation, all instruction halted until next key event).
func (h *Emulator) OpcodeFX0A(op chip8.WORD) {
	if debug {
		fmt.Printf("FX0A: %x\n", op)
	}
	fmt.Println("FX0A Not implemented!")
}

// Sets the delay timer to VX.
func (h *Emulator) OpcodeFX15(op chip8.WORD) {
	if debug {
		fmt.Printf("FX15: %x\n", op)
	}
	regx, _ := chip8.GetXYReg(op)
	h.Delay = h.Registers[regx]
}

// Sets the sound timer to VX.
func (h *Emulator) OpcodeFX18(op chip8.WORD) {
	if debug {
		fmt.Printf("FX18: %x\n", op)
	}
	fmt.Println("FX18 Not implemented!")
}

// Adds VX to I. VF is not affected
func (h *Emulator) OpcodeFX1E(op chip8.WORD) {
	if debug {
		fmt.Printf("FX1E: %x\n", op)
	}
	regx := (op & 0x0F00) >> 8
	if h.I+chip8.WORD(h.Registers[regx]) > 0xFFF {
		h.Registers[0xF] = 1
	} else {
		h.Registers[0xF] = 0
	}
	h.I += chip8.WORD(h.Registers[regx])
}

// Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
func (h *Emulator) OpcodeFX29(op chip8.WORD) {
	if debug {
		fmt.Printf("FX29: %x\n", op)
	}
	regx := (op & 0x0F00) >> 8
	h.I = chip8.WORD(h.Registers[regx]) * 0x5
}

// Stores the binary-coded decimal representation of VX, with the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
func (h *Emulator) OpcodeFX33(op chip8.WORD) {
	if debug {
		fmt.Printf("FX33: %x\n", op)
	}
	regx := (op & 0x0F00) >> 8
	value := h.Registers[regx]

	hundreds := value / 100
	tens := (value / 10) % 10
	units := value % 10

	h.Memory[h.I] = hundreds
	h.Memory[h.I+1] = tens
	h.Memory[h.I+2] = units
}

// Stores from V0 to VX (including VX) in memory, starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.[d]
func (h *Emulator) OpcodeFX55(op chip8.WORD) {
	if debug {
		fmt.Printf("FX55: %x\n", op)
	}
	regx := (op & 0x0F00) >> 8
	for i := 0; chip8.WORD(i) <= regx; i++ {
		h.Memory[h.I+chip8.WORD(i)] = h.Registers[i]
	}
	h.I = h.I + regx + 1
}

// Fills from V0 to VX (including VX) with values from memory, starting at address I. The offset from I is increased by 1 for each value read, but I itself is left unmodified.[d]
func (h *Emulator) OpcodeFX65(op chip8.WORD) {
	if debug {
		fmt.Printf("FX65: %x\n", op)
	}
	regx := (op & 0x0F00) >> 8

	for i := 0; i <= int(regx); i++ {
		h.Registers[i] = h.Memory[int(h.I)+i]
	}
}
