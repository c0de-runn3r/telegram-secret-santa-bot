package telegram

import (
	"main/clients/telegram"
)

var (
	ButtonCreateGame    = NewButton("Створити нову гру")
	ButtonConnectToGame = NewButton("Приєднатись до гри")
	ButtonMain          = NewButton("На головну")
)

var (
	ActionKeyboard = telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{*ButtonCreateGame, *ButtonConnectToGame},
			{*ButtonMain},
		},
		ResizeKeyboard: true,
	}
)

func NewButton(text string) *telegram.KeyboardButton {
	return &telegram.KeyboardButton{
		Text: text,
	}
}
