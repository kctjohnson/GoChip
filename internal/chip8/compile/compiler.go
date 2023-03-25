package compile

import "chip8-emu/internal/chip8/parse"

type Compiler struct {
	tokens []parse.Token
}

func NewCompiler(tokens []parse.Token) *Compiler {
	return &Compiler{tokens: tokens}
}
