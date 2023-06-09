package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kctjohnson/chip8-emu/internal/chip8/compiler"
	"github.com/kctjohnson/chip8-emu/internal/chip8/parser"
)

func main() {
	inputPath := flag.String("in", "", "Input file")
	outputPath := flag.String("out", "", "Output file")
	flag.Parse()

	if *inputPath == "" {
		fmt.Printf("Missing input path argument")
		return
	}

	if *outputPath == "" {
		fmt.Printf("Missing output path argument")
		return
	}

	file, err := os.ReadFile(*inputPath)
	if err != nil {
		panic(err)
	}

	p := parser.NewParser(string(file))
	p.ReadTokens()

	c := compiler.NewCompiler(p.GetTokens())
	data := c.Compile()
	err = os.WriteFile(*outputPath, data, 0777)
	if err != nil {
		panic(err)
	}
}
