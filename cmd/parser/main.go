package main

import (
	"chip8-emu/internal/chip8/compile"
	"chip8-emu/internal/chip8/parse"
	"os"
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
