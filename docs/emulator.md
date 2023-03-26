# Chip-8 Emulator

## Description

Yet another chip-8 emulator. There are better examples and emulators out there,
so I wouldn't recommend using this one if you plan on actually playing chip-8 games.

## Running

The emulator takes in any Chip-8 rom file, either ones you've found, or written and compiled.  

`go run ./cmd/tui -in FILE_PATH`  

## Keybindings

The left side of the keyboard is mapped to the chip-8 keys.  

```
1 2 3 4
Q W E R
A S D F
Z X C V
```

### Other Keybindings

`p: Pause the game`  
`n: When paused, step one instruction forward`  
`?: Switch to debug mode`  
