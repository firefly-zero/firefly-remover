package main

import (
	"strings"

	"github.com/firefly-zero/firefly-go/firefly"
)

type State struct {
	settings firefly.Settings
	authorID string
	appID    string
	font     firefly.Font
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
	if target.Raw != nil {
		authorID, appID, _ := strings.Cut(string(target.Raw), ".")
		state.authorID = authorID
		state.appID = appID
	}
}

func update() {
	// ...
}

func render() {
	firefly.ClearScreen(firefly.ColorWhite)
	drawBackgroundGrid()
	drawBackgroundBox()
}

func drawBackgroundGrid() {
	const cellSize = 10
	firefly.ClearScreen(firefly.ColorWhite)
	lineStyle := firefly.L(firefly.ColorLightGray, 1)
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
	size := firefly.S(firefly.Width-margin*2, firefly.Height-margin*2)
	firefly.DrawRoundedRect(
		firefly.P(margin+1, margin+1),
		size,
		firefly.S(4, 4),
		firefly.Solid(firefly.ColorBlack),
	)
	firefly.DrawRoundedRect(
		firefly.P(margin, margin),
		size,
		firefly.S(4, 4),
		firefly.Style{
			FillColor:   firefly.ColorWhite,
			StrokeColor: firefly.ColorBlack,
			StrokeWidth: 1,
		},
	)
}
