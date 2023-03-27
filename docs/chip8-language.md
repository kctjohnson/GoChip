# Chip-8 Language

## Intro

This was entirely made up based on what I saw on the wiki, so if it differes from the output
of other chip-8 disassemblers, it's because I tossed it together so that I could read the
rom code easier.

## Example

```CLS
MOV REG[0x0], 0
ADD REG[0x0], 1
MOV REG[0x1], REG[0x0]
SNEQ REG[0x0], 0xF
JMP 0x202
JMP 0x204
```

## Commands

The format column is how the instruction is formatted.  

```CMD
CMD_VAL
CMD_REG
CMD_REG_VAL
CMD_REG_REG
CMD_REG_SPC
CMD_SPC_REG
CMD_SPC_VAL
```  

All addresses must be provided starting at 0x200.  


| Command | Description                                                                         | Format          | Example                     |
| ------- | ----------------------------------------------------------------------------------- | --------------- | --------------------------- |
| CLS     | Clears the screen                                                                   | CMD             | `CLS`                       |
| SYSCALL | Makes a system call                                                                 | CMD_VAL         | `SYSCALL 0xNNN`             |
| CALL    | Makes a call to the given subroutine                                                | CMD_VAL         | `CALL 0xNNN`                |
| RET     | Returns from a subroutine                                                           | CMD             | `RET`                       |
| JMP     | Jumps to the given address                                                          | CMD_VAL         | `JMP 0xNNN`                 |
| RJMP    | Jumps to the address plus register 0                                                | CMD_VAL         | `RJMP 0xNNN`                |
| SEQ     | Skips the next instruction if equal                                                 | CMD_REG_REG     | `SEQ REG[0xN], REG[0xN]`    |
| SEQ     | Skips the next instruction if equal                                                 | CMD_REG_VAL     | `SEQ REG[0xN], 0xNN`        |
| SNEQ    | Skips the next instruction if not equal                                             | CMD_REG_REG     | `SNEQ REG[0xN], REG[0xN]`   |
| SNEQ    | Skips the next instruction if not equal                                             | CMD_REG_VAL     | `SNEQ REG[0xN], 0xNN`       |
| JKP     | Skip the next instruction if the key in the register is pressed                     | CMD_REG         | `JKP REG[0xN]`              |
| JKNP    | Skip the next instruction if the key in the register is not pressed                 | CMD_REG         | `JKNP REG[0xN]`             |
| WK      | Waits for the key in the register to be pressed                                     | CMD_REG         | `WK REG[0xN]`               |
| MOV     | Move the value from the right register to the left                                  | CMD_REG_REG     | `MOV REG[0xN], REG[0xN]`    |
| MOV     | Move the time delay into the register                                               | CMD_REG_SPC     | `MOV REG[0xN], DELAY`       |
| MOV     | Move the value into the register                                                    | CMD_REG_VAL     | `MOV REG[0xN], 0xNN`        |
| MOV     | Move the given value into the ADP memory pointer                                    | CMD_SPC_VAL     | `MOV ADP, 0xNN`             |
| MOV     | Move the given register value into the timer delay                                  | CMD_SPC_REG     | `MOV DELAY, REG[0xN]`       |
| MOV     | Move the given register value into the sound timer delay                            | CMD_SPC_REG     | `MOV SND_DELAY, REG[0xN]`   |
| ADD     | Add the right register to the left, store the value in the left                     | CMD_REG_REG     | `ADD REG[0xN], REG[0xN]`    |
| ADD     | Add the value to the register                                                       | CMD_REG_VAL     | `ADD REG[0xN], 0xNN`        |
| ADD     | Add the register to the ADP memory pointer                                          | CMD_SPC_REG     | `ADD ADP, REG[0xN]`         |
| SUB     | Subtract the right register from the left, store the value in the left              | CMD_REG_REG     | `SUB REG[0xN], REG[0xN]`    |
| OR      | Bitwise OR                                                                          | CMD_REG_REG     | `OR REG[0xN], REG[0xN]`     |
| AND     | Bitwise AND                                                                         | CMD_REG_REG     | `AND REG[0xN], REG[0xN]`    |
| XOR     | Bitwise XOR                                                                         | CMD_REG_REG     | `XOR REG[0xN], REG[0xN]`    |
| SHR     | Shift right one bit                                                                 | CMD_REG         | `SHR REG[0xN]`              |
| SHL     | Shift left one bit                                                                  | CMD_REG         | `SHL REG[0xN]`              |
| BRND    | Generates a random number, and bitwise ANDs it with the given register              | CMD_REG         | `BRND REG[0xN]`             |
| DRW     | Draws the sprite at the ADP memory pointer at the given register X and Y            | CMD_REG_REG_VAL | `DRW REG[0xN], REG[0xN], N` |
| FX29    | Sets ADP to the location of the sprite in the given register                        | CMD_REG         | `FX29 REG[0xN]`             |
| FX33    | Stores the binary-coded decimal representation of VX (Check Wiki for more info)     | CMD_REG         | `FX33 REG[0xN]`             |
| FX55    | Stores from V0 to VX (include VX) in memory, starting at ADP                        | CMD_REG         | `FX55 REG[0xN]`             |
| FX65    | Fills from V0 to VX (including VX) with values from memory, starting at address ADP | CMD_REG         | `FX65 REG[0xN]`             |
