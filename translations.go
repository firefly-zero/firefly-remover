package main

import "github.com/firefly-zero/firefly-go/firefly"

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

func msgAlreadyRemoved() string {
	switch state.settings.Language {
	case firefly.Dutch:
		return "app already removed" // TODO
	case firefly.French:
		return "app already removed" // TODO
	case firefly.German:
		return "app already removed" // TODO
	case firefly.Italian:
		return "app already removed" // TODO
	case firefly.Polish:
		return "app already removed" // TODO
	case firefly.Russian:
		return "приложение уже удалено"
	case firefly.Spanish:
		return "app already removed" // TODO
	case firefly.Swedish:
		return "app already removed" // TODO
	case firefly.TokiPona:
		return "app already removed" // TODO
	case firefly.Turkish:
		return "app already removed" // TODO
	case firefly.Ukrainian:
		return "app already removed" // TODO
	}
	return "app already removed"
}

func msgROM() string {
	switch state.settings.Language {
	case firefly.Dutch:
		return "app ROM" // TODO
	case firefly.French:
		return "app ROM" // TODO
	case firefly.German:
		return "app ROM" // TODO
	case firefly.Italian:
		return "app ROM" // TODO
	case firefly.Polish:
		return "app ROM" // TODO
	case firefly.Russian:
		return "app ROM"
	case firefly.Spanish:
		return "app ROM" // TODO
	case firefly.Swedish:
		return "app ROM" // TODO
	case firefly.TokiPona:
		return "app ROM" // TODO
	case firefly.Turkish:
		return "app ROM" // TODO
	case firefly.Ukrainian:
		return "app ROM" // TODO
	}
	return "app ROM"
}

func msgData() string {
	switch state.settings.Language {
	case firefly.Dutch:
		return "data and save files" // TODO
	case firefly.French:
		return "data and save files" // TODO
	case firefly.German:
		return "data and save files" // TODO
	case firefly.Italian:
		return "data and save files" // TODO
	case firefly.Polish:
		return "data and save files" // TODO
	case firefly.Russian:
		return "данные"
	case firefly.Spanish:
		return "data and save files" // TODO
	case firefly.Swedish:
		return "data and save files" // TODO
	case firefly.TokiPona:
		return "data and save files" // TODO
	case firefly.Turkish:
		return "data and save files" // TODO
	case firefly.Ukrainian:
		return "data and save files" // TODO
	}
	return "data and save files"
}

func msgShots() string {
	switch state.settings.Language {
	case firefly.Dutch:
		return "screenshots" // TODO
	case firefly.French:
		return "screenshots" // TODO
	case firefly.German:
		return "screenshots" // TODO
	case firefly.Italian:
		return "screenshots" // TODO
	case firefly.Polish:
		return "screenshots" // TODO
	case firefly.Russian:
		return "скриншоты"
	case firefly.Spanish:
		return "screenshots" // TODO
	case firefly.Swedish:
		return "screenshots" // TODO
	case firefly.TokiPona:
		return "screenshots" // TODO
	case firefly.Turkish:
		return "screenshots" // TODO
	case firefly.Ukrainian:
		return "screenshots" // TODO
	}
	return "screenshots"
}
