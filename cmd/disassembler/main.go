package main

import (
	"flag"
	"fmt"

	"github.com/kctjohnson/chip8-emu/internal/chip8/disassembler"
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
