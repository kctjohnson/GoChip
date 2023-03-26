package main

import (
	"chip8-emu/internal/chip8/compile"
	"chip8-emu/internal/chip8/parse"
	"flag"
	"fmt"
	"os"
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

	p := parse.NewParser(string(file))
	p.ReadTokens()

	c := compile.NewCompiler(p.GetTokens())
	data := c.Compile()
	err = os.WriteFile(*outputPath, data, 0777)
	if err != nil {
		panic(err)
	}
}
