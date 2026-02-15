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

type Line struct {
	kind     Kind
	selected bool
}

type State struct {
	settings firefly.Settings
	authorID string
	appID    string
	font     firefly.Font
	lines    []Line
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
		state.lines = append(state.lines, Line{kind: KindROM})
	}
	if firefly.FileExists("data/" + authorID + "/" + appID + "/stats") {
		state.lines = append(state.lines, Line{kind: KindData})
	}
	if firefly.FileExists("data/" + authorID + "/" + appID + "/shots/001.ffs") {
		state.lines = append(state.lines, Line{kind: KindShots})
	}
}

func update() {
	// ...
}

func render() {
	drawBackgroundGrid()
	drawBackgroundBox()
	if state.authorID == "" {
		drawCentered(msgNoTarget())
		return
	}
	if len(state.lines) == 0 {
		drawCentered(msgAlreadyRemoved())
		return
	}

	drawHeader(1, "What do you want to delete?")
	for i, line := range state.lines {
		drawLine(i, line)
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

func drawLine(i int, line Line) {
	text := lineMsg(line.kind)
	point := firefly.P(30, 20+state.font.CharHeight()*(i+2))
	firefly.DrawText(
		text,
		state.font,
		point,
		state.settings.Theme.Primary,
	)
}

func lineMsg(kind Kind) string {
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
