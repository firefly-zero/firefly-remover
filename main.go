package main

import (
	"strings"

	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/firefly-zero/firefly-go/firefly/sudo"
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

	targetFile := firefly.LoadFile("target", nil)
	if targetFile.Raw == nil {
		return
	}
	target := strings.Trim(string(targetFile.Raw), " \n")
	authorID, appID, _ := strings.Cut(target, ".")
	state.authorID = authorID
	state.appID = appID

	if sudo.FileExists("roms/" + authorID + "/" + appID + "/_bin") {
		state.options = append(state.options, Options{kind: KindROM})
	}
	if hasData(authorID, appID) {
		state.options = append(state.options, Options{kind: KindData})
	}
	if sudo.FileExists("data/" + authorID + "/" + appID + "/shots/001.ffs") {
		state.options = append(state.options, Options{kind: KindShots})
	}
}

func hasData(authorID, appID string) bool {
	if sudo.FileExists("data/" + authorID + "/" + appID + "/stash") {
		return true
	}
	dataFiles := sudo.ListFiles("data/" + authorID + "/" + appID + "/etc")
	return len(dataFiles) != 0
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
			option := &state.options[state.cursor]
			option.selected = !option.selected
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

	drawHeader("What do you want to delete?")
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

func drawHeader(text string) {
	const margin = 25
	firefly.DrawText(
		text,
		state.font,
		firefly.P(margin, margin+state.font.CharHeight()),
		state.settings.Theme.Accent,
	)
}

func drawOption(i int, option Options) {
	const margin = 25
	h := state.font.CharHeight()

	theme := state.settings.Theme
	text := msgOption(option.kind)
	lineH := state.font.CharHeight() + 2
	point := firefly.P(margin, margin+lineH*(i+2))
	firefly.DrawText(
		text,
		state.font,
		point,
		theme.Primary,
	)

	point.X = firefly.Width - margin - h*2
	point.Y -= h - 2
	{
		switchPoint := point
		color := theme.Secondary
		if option.selected {
			switchPoint.X += h
			color = theme.Accent
		}
		firefly.DrawCircle(switchPoint, h, firefly.Solid(color))
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
