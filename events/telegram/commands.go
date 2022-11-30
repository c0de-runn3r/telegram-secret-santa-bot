package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"main/clients/telegram"
	storage "main/files_storage"
	. "main/fsm"
)

func (p *Processor) doMessage(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%s'", text, username)
	state := FSM.CurrentState
	switch text {
	case StartCmd:
		FSM.SetState(*ActionState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgHello, Keyboard: &ActionKeyboard})
	case HelpCmd:
		FSM.SetState(*ActionState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgHelp})
	case cmdMain:
		FSM.SetState(*ActionState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgCancel})
	default:
		switch state {
		case *ActionState:
			p.ProcessAction(text, chatID, username)
		case *NewGameNameState:
			FSM.SetState(*ActionState)
			p.CreateNewGame(text, chatID, username)
		case *ConnectExistingGameState:
			FSM.SetState(*GameSettingsState)
			p.ConnectToExistingGame(text, chatID, username)
		case *GameSettingsState:
			p.ChangeGameSettings(text, chatID, username)
		case *UpdateWishesState:
			p.UpdateWishes(text, chatID, username)
			FSM.SetState(*ActionState)
		default:
			FSM.SetState(*ActionState)
			p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgUnknownCommand})
		}
	}

	return nil
}

func (p *Processor) ProcessAction(text string, chatID int, username string) {
	switch text {
	case cmdCreateNewGame:
		FSM.SetState(*NewGameNameState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgNameNewGame})
	case cmdConnectToExistingGame:
		FSM.SetState(*ConnectExistingGameState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgSendIDOfGame})
	case cmdCheckMyGames:
		p.CheckGames(text, chatID, username)
	default:
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgUnknownCommand})
	}

}

func (p *Processor) CreateNewGame(gameName string, chatID int, username string) {
	log.Printf("creating new game [%s]", gameName)
	id := storage.DB.AddNewGame(gameName, username)
	msg := fmt.Sprintf("Нову гру %s створено. ID: %v", gameName, id)
	p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msg})
}

func (p *Processor) ConnectToExistingGame(strID string, chatID int, username string) {
	gameID, err := strconv.Atoi(strID)
	if err != nil {
		log.Println("Can't convert stringID into int")
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgSendIntNotStr})
		return
	}
	var game storage.Game
	storage.DB.First(&game, gameID)
	if game.ID != 0 {
		msg := fmt.Sprintf("Вітаю!\nВи приєднались до %s", game.Name)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msg})
		storage.DB.AddUserToGame(&game, username)
	} else {
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgUndefinedGameID})
		FSM.SetState(*ConnectExistingGameState)
	}
}

func (p *Processor) CheckGames(text string, chatID int, username string) {
	msg := "Твої ігри:"
	var games []*storage.SantaUser
	storage.DB.Table("santa_users").Where("username = ?", username).Find(&games)
	for i := 0; i < len(games); i++ {
		if games[i].Game != "" {
			msg = fmt.Sprintf("%s\n%s", msg, games[i].Game)
		}
	}
	p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msg})
}

func (p *Processor) ChangeGameSettings(text string, chatID int, username string) {
	switch text {
	case cmdChangeWishes:
		FSM.SetState(*UpdateWishesState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgAddWishes})
	default:
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgUnknownCommand})
	}
}

func (p *Processor) UpdateWishes(text string, chatID int, username string) {
	// check to what game add wishes to
	storage.DB.AddOrUpdateWishes(text, username)
	p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgWishesAdded})
}
