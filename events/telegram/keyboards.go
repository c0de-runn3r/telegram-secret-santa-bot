package telegram

import (
	"main/clients/telegram"
)

var (
	ButtonCreateGame    = NewButton("Створити нову гру")
	ButtonConnectToGame = NewButton("Приєднатись до гри")
	ButtonCheckMyGames  = NewButton("Переглянути мої ігри")
	ButtonMain          = NewButton("На головну")
)

var (
	ActionKeyboard = telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{*ButtonCreateGame, *ButtonConnectToGame},
			{*ButtonCheckMyGames},
			{*ButtonMain},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
)

func NewButton(text string) *telegram.KeyboardButton {
	return &telegram.KeyboardButton{
		Text: text,
	}
}

func AddButtonToKeyboard(button *telegram.KeyboardButton, keyboard *telegram.ReplyKeyboardMarkup, n int) {
	keyboard.Keyboard[n] = append(keyboard.Keyboard[n], *button)
}
