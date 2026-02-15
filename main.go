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
	drawBackgroundGrid()
	drawBackgroundBox()
	if state.authorID == "" {
		drawNoTarget()
		return
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

func drawNoTarget() {
	text := msgNoTarget()
	x := (firefly.Width - state.font.LineWidth(text)) / 2
	y := (firefly.Height - state.font.CharHeight()) / 2
	firefly.DrawText(
		text,
		state.font,
		firefly.P(x, y),
		state.settings.Theme.Primary,
	)
}

func msgNoTarget() string {
	switch state.settings.Language {
	case firefly.Dutch:
		return "no app selected" // TODO
	case firefly.French:
		return "no app selected" // TODO
	case firefly.German:
		return "no app selected" // TODO
	case firefly.Italian:
		return "no app selected" // TODO
	case firefly.Polish:
		return "no app selected" // TODO
	case firefly.Russian:
		return "приложение не выбрано"
	case firefly.Spanish:
		return "no app selected" // TODO
	case firefly.Swedish:
		return "no app selected" // TODO
	case firefly.TokiPona:
		return "no app selected" // TODO
	case firefly.Turkish:
		return "no app selected" // TODO
	case firefly.Ukrainian:
		return "no app selected" // TODO
	}
	return "no app selected"
}
