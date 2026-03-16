package main

import (
	"strings"

	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/firefly-zero/firefly-go/firefly/sudo"
)

const (
	// For how long to show messages on the screen before exiting.
	defaultMsgTTL = 4 * 60
	// The distance between the content text and screen borders.
	contentMargin = 20
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
	msg      string
	msgTTL   int
	cursor   int
	dpad     firefly.DPad4
	btns     firefly.Buttons
	peer     firefly.Peer
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
	if targetFile == nil {
		state.msg = msgNoTarget()
		state.msgTTL = defaultMsgTTL
		return
	}
	target := strings.Trim(string(targetFile), " \n")
	authorID, appID, _ := strings.Cut(target, ".")
	state.authorID = authorID
	state.appID = appID

	if sudo.FileExists("roms/" + authorID + "/" + appID + "/_bin") {
		state.options = append(state.options, Options{kind: KindROM})
	}
	if hasData(authorID, appID) {
		state.options = append(state.options, Options{kind: KindData})
	}
	if hasShots(authorID, appID) {
		state.options = append(state.options, Options{kind: KindShots})
	}
	if len(state.options) == 0 {
		state.msg = msgAlreadyRemoved()
		state.msgTTL = defaultMsgTTL
	}

	// Set the peer to the current device.
	// While most games must avoid state drift, in case of the remover app
	// we DO want every device to run their own copy of the app,
	// even in multiplayer. It shouldn't be possible from one device
	// to remove apps and data from another device.
	for _, peer := range firefly.GetPeers().Slice() {
		if me.Eq(peer) {
			state.peer = peer
		}
	}
}

func hasData(authorID, appID string) bool {
	if sudo.FileExists("data/" + authorID + "/" + appID + "/stash") {
		return true
	}
	dataFiles := sudo.ListFiles("data/" + authorID + "/" + appID + "/etc")
	return len(dataFiles) != 0
}

func hasShots(authorID, appID string) bool {
	dataFiles := sudo.ListFiles("data/" + authorID + "/" + appID + "/shots")
	return len(dataFiles) != 0
}

func update() {
	// Automatically exit if there is a message on the screen
	// and we've been showing it for long enough.
	if state.msg != "" {
		if state.msgTTL == 0 {
			firefly.Quit()
			return
		}
		state.msgTTL--
	}

	handlePad()
	handleButtons()
}

func handlePad() {
	pad, _ := firefly.ReadPad(state.peer)
	dpad := pad.DPad4()
	pressed := dpad.JustPressed(state.dpad)
	state.dpad = dpad
	switch pressed {
	case firefly.DPad4Down:
		if state.cursor < len(state.options) {
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
	btns := firefly.ReadButtons(state.peer)
	released := btns.JustReleased(state.btns)
	state.btns = btns
	if released.W {
		firefly.Quit()
		return
	}
	if released.S || released.E {
		if state.msg != "" {
			firefly.Quit()
			return
		}
		if state.cursor < len(state.options) {
			option := &state.options[state.cursor]
			option.selected = !option.selected
		} else {
			if anySelected() {
				removeApp()
			} else {
				firefly.Quit()
			}
		}
	}
}

func anySelected() bool {
	for _, option := range state.options {
		if option.selected {
			return true
		}
	}
	return false
}

func removeApp() {
	authorID := state.authorID
	appID := state.appID

	// Delete ROM and detect if we can remove the whole data dir
	// or only some subdirs.
	delAllData := true
	delEtc := false
	delShots := false
	for _, option := range state.options {
		if option.kind == KindROM {
			if option.selected {
				sudo.RemoveDir("roms/" + authorID + "/" + appID)
				// Reset launcher cache.
				sudo.RemoveFile("data/sys/launcher/etc/metas")
			}
			continue
		}
		if !option.selected {
			delAllData = false
			continue
		}
		if option.kind == KindData {
			delEtc = true
		}
		if option.kind == KindShots {
			delShots = true
		}
	}

	state.msg = msgRemoved()
	state.msgTTL = defaultMsgTTL

	if delAllData {
		sudo.RemoveDir("data/" + authorID + "/" + appID)
		return
	}

	// Delete the app data and stash.
	if delEtc {
		sudo.RemoveDir("data/" + authorID + "/" + appID + "/etc")
		stashPath := "data/" + authorID + "/" + appID + "/stash"
		if sudo.FileExists(stashPath) {
			sudo.RemoveFile(stashPath)
		}
	}

	// Delete screenshots.
	if delShots {
		sudo.RemoveDir("data/" + authorID + "/" + appID + "/shots")
	}
}

func render() {
	drawBackgroundGrid()
	drawBackgroundBox()
	if state.msg != "" {
		drawCentered(state.msg)
		return
	}

	drawCursor()
	drawHeader("What do you want to delete?")
	for i, option := range state.options {
		drawOption(i, option)
	}
	drawButton()
}

func drawBackgroundGrid() {
	const cellSize = 8

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
	const margin = 12

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
	firefly.DrawText(
		text,
		state.font,
		firefly.P(contentMargin, contentMargin+state.font.CharHeight()-4),
		state.settings.Theme.Accent,
	)
}

func drawOption(i int, option Options) {
	h := state.font.CharHeight()
	theme := state.settings.Theme
	text := msgOption(option.kind)
	lineH := state.font.CharHeight() + 4
	point := firefly.P(contentMargin, contentMargin+lineH*(i+2)-8)
	firefly.DrawText(
		text,
		state.font,
		point,
		theme.Primary,
	)

	point.X = firefly.Width - contentMargin - h*2
	point.Y -= h - 3
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

func drawButton() {
	var msg string
	if anySelected() {
		msg = msgRemove()
	} else {
		msg = msgCancel()
	}
	lineH := state.font.CharHeight() + 4
	lineNo := len(state.options) + 2
	firefly.DrawText(
		msg,
		state.font,
		firefly.P(contentMargin, contentMargin+lineH*lineNo-8),
		state.settings.Theme.Accent,
	)
}

func drawCursor() {
	const margin = contentMargin - 4
	theme := state.settings.Theme
	lineH := state.font.CharHeight() + 4
	point := firefly.P(margin, contentMargin+lineH*(state.cursor+1)-3)
	size := firefly.S(firefly.Width-margin*2, lineH)
	corners := firefly.S(4, 4)

	{
		point := firefly.P(point.X+1, point.Y+1)
		firefly.DrawRoundedRect(point, size, corners, firefly.Solid(theme.Primary))
	}

	style := firefly.Style{
		FillColor:   theme.BG,
		StrokeColor: theme.Primary,
		StrokeWidth: 1,
	}
	firefly.DrawRoundedRect(point, size, corners, style)
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
