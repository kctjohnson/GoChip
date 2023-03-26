# Disassembler

## Description

Disassembles Chip-8 opcodes into a readable assembly-style language.  

## Running

The disassembler takes in any Chip-8 rom file, either ones you've found, or written and compiled.  
  
```❯ go run ./cmd/disassembler -in compiled.o
0x00E0  CLS
0x6000  MOV reg[0x0], 0
0x7001  ADD reg[0x0], 1
0x8100  MOV reg[0x1], reg[0x0]
0x400F  SNEQ reg[0x0], 15
0x1202  JMP 0x0202
0x1204  JMP 0x0204
0x0000  SYSCALL 0x0000
0x0000  SYSCALL 0x0000
0x0000  SYSCALL 0x0000
0x0000  SYSCALL 0x0000
0x0000  SYSCALL 0x0000
```
