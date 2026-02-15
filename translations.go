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
