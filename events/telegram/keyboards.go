package telegram

import (
	"fmt"
	"main/clients/telegram"
	storage "main/files_storage"
	"strconv"
)

var (
	ButtonCreateGame   = NewButton(cmdCreateNewGame)
	ButtonCheckMyGames = NewButton(cmdCheckMyGames)
	ButtonMain         = NewButton(cmdMain)
)

func makeActionKeyboard(username string) *telegram.ReplyKeyboardMarkup {
	var games []*storage.SantaUser
	storage.DB.Table("santa_users").Where("username = ?", username).Find(&games)
	keyboard := &telegram.ReplyKeyboardMarkup{
		Keyboard:        make([][]telegram.KeyboardButton, len(games)+2),
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	for i := 0; i < len(games); i++ {
		if games[i].Game != "" {
			buttonName := fmt.Sprintf("Налаштування гри: %s (ID:%v)", games[i].Game, games[i].SantaID)
			button := NewButton(buttonName)
			AddButtonToKeyboard(button, keyboard, i+1)
		}
	}
	AddButtonToKeyboard(ButtonCreateGame, keyboard, 0)
	AddButtonToKeyboard(ButtonMain, keyboard, len(games)+1)
	return keyboard
}

func NewButton(text string) *telegram.KeyboardButton {
	return &telegram.KeyboardButton{
		Text: text,
	}
}

func AddButtonToKeyboard(button *telegram.KeyboardButton, keyboard *telegram.ReplyKeyboardMarkup, n int) {
	keyboard.Keyboard[n] = append(keyboard.Keyboard[n], *button)
}

func createSettingsKeyboard(text string, username string, id string) *telegram.InlineKeyboardMarkup {

	idInt, _ := strconv.Atoi(id)

	showAllPlayersButton := &telegram.InlineKeyboardButton{
		Text:         cmdShowAllPlayers,
		CallbackData: "all_players " + id,
	}
	addWishesButton := &telegram.InlineKeyboardButton{
		Text:         cmdChangeWishes,
		CallbackData: "change_wishes " + id,
	}
	startGameButton := &telegram.InlineKeyboardButton{
		Text:         cmdStartGame,
		CallbackData: "start_game " + id,
	}
	changeBudgetButton := &telegram.InlineKeyboardButton{
		Text:         cmdChangeBudget,
		CallbackData: "change_budget " + id,
	}
	quitGameButton := &telegram.InlineKeyboardButton{
		Text:         cmdQuitGame,
		CallbackData: "quit_game " + id,
	}
	deleteGameButton := &telegram.InlineKeyboardButton{
		Text:         cmdDeleteGame,
		CallbackData: "quit_game " + id,
	}

	keyboard := &telegram.InlineKeyboardMarkup{} // TODO make creating this keyboard is separate func so I can pass it to user when creating game
	admin, _ := storage.DB.QueryAdmin(idInt)
	if username != admin {
		keyboard = &telegram.InlineKeyboardMarkup{
			Buttons: [][]telegram.InlineKeyboardButton{{*showAllPlayersButton}, {*addWishesButton}, {*quitGameButton}},
		}
	}
	if username == admin {
		keyboard = &telegram.InlineKeyboardMarkup{
			Buttons: [][]telegram.InlineKeyboardButton{{*showAllPlayersButton}, {*addWishesButton}, {*changeBudgetButton}, {*startGameButton}, {*deleteGameButton}},
		}
	}
	return keyboard
}
