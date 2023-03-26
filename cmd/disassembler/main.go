package main

import (
	"chip8-emu/internal/chip8/disassembler"
	"flag"
	"fmt"
)

func main() {
	inputPath := flag.String("in", "", "Input file")
	flag.Parse()

	if *inputPath == "" {
		fmt.Printf("Missing input path argument")
		return
	}

	disassembler.Disassemble(*inputPath)
}
