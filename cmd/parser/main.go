package main

import (
	"os"

	"github.com/kctjohnson/chip8-emu/internal/chip8/compile"
	"github.com/kctjohnson/chip8-emu/internal/chip8/parse"
)

func main() {
	file, err := os.ReadFile("DisassembledBrix.asm")
	if err != nil {
		panic(err)
	}
	p := parse.NewParser(string(file))
	p.ReadTokens()

	c := compile.NewCompiler(p.GetTokens())
	data := c.Compile()
	err = os.WriteFile("RecompiledBrix.chp8", data, 0777)
	if err != nil {
		panic(err)
	}
}
