package main

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kctjohnson/chip8-emu/internal/chip8/disassembler"
	"github.com/kctjohnson/chip8-emu/internal/chip8/emulator"
)

var (
	speed        = true
	displayDebug = false
)

const FPS = 240

type TickMsg time.Time

type Model struct {
	emu *emulator.Emulator
}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m Model) tick() tea.Cmd {
	fps := FPS
	if !speed {
		fps = 1
	}
	return tea.Tick(time.Second/time.Duration(fps), func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if speed {
			m.emu.Step()
		}
		return m, m.tick()
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.emu.Inputs[0] = 1
		case "2":
			m.emu.Inputs[1] = 1
		case "3":
			m.emu.Inputs[2] = 1
		case "4":
			m.emu.Inputs[3] = 1
		case "q":
			m.emu.Inputs[4] = 1
		case "w":
			m.emu.Inputs[5] = 1
		case "e":
			m.emu.Inputs[6] = 1
		case "r":
			m.emu.Inputs[7] = 1
		case "a":
			m.emu.Inputs[8] = 1
		case "s":
			m.emu.Inputs[9] = 1
		case "d":
			m.emu.Inputs[10] = 1
		case "f":
			m.emu.Inputs[11] = 1
		case "z":
			m.emu.Inputs[12] = 1
		case "x":
			m.emu.Inputs[13] = 1
		case "c":
			m.emu.Inputs[14] = 1
		case "v":
			m.emu.Inputs[15] = 1
		case "m":
			m.emu.DelayEnabled = !m.emu.DelayEnabled
		case "p":
			speed = !speed
		case "ctrl+r":
			m.emu.CPUReset()
		case "?":
			displayDebug = !displayDebug
			m.emu.LastTick = time.Now()
		case "n":
			m.emu.Step()
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if displayDebug {
		return m.debugView()
	}
	return m.gameView()
}

func (m Model) debugView() string {
	debugData := fmt.Sprintf("CurOp: %X\nDelay: %d\nADP: %X\nPC: %X\n", m.emu.CurrentOpcode, m.emu.Delay, m.emu.I, m.emu.PC)
	for i, r := range m.emu.Registers {
		debugData += fmt.Sprintf("Reg%X: %X\n", i, r)
	}
	debugData = lipgloss.PlaceHorizontal(25, lipgloss.Top, debugData)

	screen := ""
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if m.emu.ScreenData[x][y] == 0 {
				screen += "."
			} else {
				screen += "#"
			}
		}
		screen += "\n"
	}

	disassembly := ""
	for pc := m.emu.PC - 10; pc <= m.emu.PC+10; pc += 2 {
		op := m.emu.GetOpcode(pc)

		if pc == m.emu.PC {
			disassembly += "> "
		}
		disassembly += fmt.Sprintf("0x%04X ", pc) + disassembler.DisassembleOpcode(op) + "\n"
	}
	disassembly = lipgloss.PlaceHorizontal(50, lipgloss.Top, disassembly)

	inputs := ""
	for i, input := range m.emu.Inputs {
		inputs += fmt.Sprintf("INPUT[%d]: %d\n", i, input)
	}
	inputs = lipgloss.PlaceHorizontal(15, lipgloss.Top, inputs)

	stack := "STACK\n"
	for _, s := range m.emu.Stack {
		stack += fmt.Sprintf("0x%X\n", s)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, debugData, disassembly, inputs, screen, stack)
}

var (
	emptyStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#333388")).
			Foreground(lipgloss.Color("#444499"))
	tileStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#333388")).
			Foreground(lipgloss.Color("#FFFFFF"))
	shadowStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#333388")).
			Foreground(lipgloss.Color("#AAAAFF"))
	renderEmpty  = emptyStyle.Render(".")
	renderTile   = tileStyle.Render("#")
	renderShadow = shadowStyle.Render("*")
)

func (m Model) gameView() string {
	screen := ""
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			switch m.emu.ScreenData[x][y] {
			case 0:
				screen += renderEmpty
			case 1:
				screen += renderTile
			case 2:
				screen += renderShadow
			}
		}
		screen += "\n"
	}
	return screen
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	model := Model{
		emu: emulator.NewEmulator("./roms/brix.rom"),
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
