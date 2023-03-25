package main

import (
	"chip8-emu/internal/chip8/parse"
	"os"
)

func main() {
	file, err := os.ReadFile("brix.chp8")
	if err != nil {
		panic(err)
	}
	p := parse.NewParser(string(file))
	p.ReadTokens()
	p.PrintTokens()
}
