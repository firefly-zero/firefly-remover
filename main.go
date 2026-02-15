package main

import (
	"strings"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Kind uint8

const (
	KindROM   Kind = 1
	KindData  Kind = 2
	KindShots Kind = 3
)

type Options struct {
	kind     Kind
	selected bool
}

type State struct {
	settings firefly.Settings
	authorID string
	appID    string
	font     firefly.Font
	options  []Options
	removed  bool
	cursor   int
	dpad     firefly.DPad4
	btns     firefly.Buttons
}

var state State

func init() {
	firefly.Boot = boot
	firefly.Update = update
	firefly.Render = render
}

func boot() {
	me := firefly.GetMe()
	state.settings = firefly.GetSettings(me)
	state.font = firefly.LoadFont("font", nil)

	target := firefly.LoadFile("target", nil)
	if target.Raw == nil {
		return
	}
	authorID, appID, _ := strings.Cut(string(target.Raw), ".")
	state.authorID = authorID
	state.appID = appID

	if firefly.FileExists("roms/" + authorID + "/" + appID + "/_bin") {
		state.options = append(state.options, Options{kind: KindROM})
	}
	if firefly.FileExists("data/" + authorID + "/" + appID + "/stats") {
		state.options = append(state.options, Options{kind: KindData})
	}
	if firefly.FileExists("data/" + authorID + "/" + appID + "/shots/001.ffs") {
		state.options = append(state.options, Options{kind: KindShots})
	}
}

func update() {
	handlePad()
	handleButtons()
}

func handlePad() {
	me := firefly.GetMe()
	pad, _ := firefly.ReadPad(me)
	dpad := pad.DPad4()
	released := dpad.JustReleased(state.dpad)
	state.dpad = dpad
	switch released {
	case firefly.DPad4Down:
		if state.cursor < len(state.options)-1 {
			state.cursor++
		}
	case firefly.DPad4Up:
		if state.cursor > 0 {
			state.cursor--
		}
	case firefly.DPad4Left:
	case firefly.DPad4Right:
	}
}

func handleButtons() {
	me := firefly.GetMe()
	btns := firefly.ReadButtons(me)
	released := btns.JustReleased(state.btns)
	state.btns = btns
	if released.W {
		firefly.Quit()
		return
	}
	if released.S || released.E {
		if state.removed {
			firefly.Quit()
			return
		}
		if state.cursor < len(state.options) {
			state.options[state.cursor].selected = true
		} else {
			removeApp()
		}
	}
}

func removeApp() {
	for _, option := range state.options {
		if !option.selected {
			continue
		}
		switch option.kind {
		case KindROM:
			removeROM()
		case KindData:
			removeData()
		case KindShots:
			removeShots()
		}
	}
	state.removed = true
}

func removeROM() {
	// ...
}

func removeData() {
	// ...
}

func removeShots() {
	// ...
}

func render() {
	drawBackgroundGrid()
	drawBackgroundBox()
	if state.authorID == "" {
		drawCentered(msgNoTarget())
		return
	}
	if len(state.options) == 0 {
		drawCentered(msgAlreadyRemoved())
		return
	}
	if state.removed {
		drawCentered(msgRemoved())
		return
	}

	drawHeader(1, "What do you want to delete?")
	for i, option := range state.options {
		drawOption(i, option)
	}
}

func drawBackgroundGrid() {
	const cellSize = 10

	theme := state.settings.Theme
	firefly.ClearScreen(theme.BG)
	lineStyle := firefly.L(theme.Secondary, 1)
	for x := cellSize; x < firefly.Width; x += cellSize {
		firefly.DrawLine(
			firefly.P(x, 0),
			firefly.P(x, firefly.Height),
			lineStyle,
		)
	}
	for y := cellSize; y < firefly.Height; y += cellSize {
		firefly.DrawLine(
			firefly.P(0, y),
			firefly.P(firefly.Width, y),
			lineStyle,
		)
	}
}

func drawBackgroundBox() {
	const margin = 15

	theme := state.settings.Theme
	size := firefly.S(firefly.Width-margin*2, firefly.Height-margin*2)
	firefly.DrawRoundedRect(
		firefly.P(margin+1, margin+1),
		size,
		firefly.S(4, 4),
		firefly.Solid(theme.Primary),
	)
	firefly.DrawRoundedRect(
		firefly.P(margin, margin),
		size,
		firefly.S(4, 4),
		firefly.Style{
			FillColor:   theme.BG,
			StrokeColor: theme.Primary,
			StrokeWidth: 1,
		},
	)
}

func drawCentered(text string) {
	x := (firefly.Width - state.font.LineWidth(text)) / 2
	y := (firefly.Height - state.font.CharHeight()) / 2
	firefly.DrawText(
		text,
		state.font,
		firefly.P(x, y),
		state.settings.Theme.Primary,
	)
}

func drawHeader(line int, text string) {
	firefly.DrawText(
		text,
		state.font,
		firefly.P(20, 20+state.font.CharHeight()*line),
		state.settings.Theme.Accent,
	)
}

func drawOption(i int, option Options) {
	theme := state.settings.Theme
	text := msgOption(option.kind)
	point := firefly.P(30, 20+state.font.CharHeight()*(i+2))
	firefly.DrawText(
		text,
		state.font,
		point,
		theme.Primary,
	)

	h := state.font.CharHeight()
	point.Y -= h
	{
		switchPoint := point
		if option.selected {
			switchPoint.X += h
		}
		firefly.DrawCircle(switchPoint, h, firefly.Solid(theme.Accent))
	}

	firefly.DrawRoundedRect(
		point,
		firefly.S(h*2, h),
		firefly.S(h/2, h/2),
		firefly.Outlined(theme.Primary, 1),
	)
}

func msgOption(kind Kind) string {
	switch kind {
	case KindData:
		return msgData()
	case KindROM:
		return msgROM()
	case KindShots:
		return msgShots()
	}
	return ""
}
