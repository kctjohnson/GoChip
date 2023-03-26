package chip8

type (
	BYTE uint8
	WORD uint16
)

func GetXYReg(op WORD) (WORD, WORD) {
	regx := op & 0x0F00 // Mask off reg x
	regx = regx >> 8    // Shift x across
	regy := op & 0x00F0 // Mask off reg y
	regy = regy >> 4    // Shift y across
	return regx, regy
}
