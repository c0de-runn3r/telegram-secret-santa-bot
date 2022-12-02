package telegram

import (
	"fmt"
	"log"
	"regexp"
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
	switch text { //for commands
	case StartCmd: // /start
		FSM.SetState(*ActionState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgHello, KeyboardReply: &ActionKeyboard})
	case HelpCmd: // /help
		FSM.SetState(*ActionState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgHelp, KeyboardReply: &ActionKeyboard})
	case cmdMain: // to the main menu
		FSM.SetState(*ActionState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgCancel, KeyboardReply: &ActionKeyboard})
	default: // other cases
		switch state {
		case *ActionState: // actions menu
			p.ProcessAction(text, chatID, username)
		case *NewGameNameState: // receive name of the new game
			FSM.SetState(*ActionState)
			p.CreateNewGame(text, chatID, username)
		case *ConnectExistingGameState: // receive id to connect to game
			FSM.SetState(*ActionState)
			p.ConnectToExistingGame(text, chatID, username)
		case *MyGamesSate: // receive id to change settings of the game
			p.ChooseTheGame(text, chatID, username)
		// case *GameSettingsState:
		// 	p.ChangeGameSettings(text, chatID, username)
		case *UpdateWishesState:
			p.UpdateWishes(text, chatID, username)
			FSM.SetState(*ActionState)
		default:
			p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgUnknownCommand})
		}
	}

	return nil
}

func (p *Processor) doCallbackQuerry(text string, chatID int, username string) error {
	log.Printf("got new callback data '%s' from '%s'", text, username)
	command, id := cutTextToData(text)
	switch command {
	case "change_wishes":
		FSM.SetState(*UpdateWishesState)
		storage.ListOfWishUpdates.Wishes = append(storage.ListOfWishUpdates.Wishes, &storage.WishUpdateInfo{
			ID:       id,
			Username: username})
		p.tg.SendMessage(telegram.MessageParams{
			ChatID: chatID,
			Text:   msgAddWishes,
		})
	default:
		p.tg.SendMessage(telegram.MessageParams{
			ChatID: chatID,
			Text:   msgSmthWrong,
		})
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
		FSM.SetState(*ConnectExistingGameState)
		return
	}
	var game storage.Game
	storage.DB.First(&game, gameID)
	if game.ID != 0 {
		msg := fmt.Sprintf("Вітаю!\nВи приєднались до %s", game.Name)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msg, KeyboardReply: &ActionKeyboard})
		storage.DB.AddUserToGame(&game, username)
	} else {
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgUndefinedGameID})
		FSM.SetState(*ConnectExistingGameState)
	}
}

func (p *Processor) CheckGames(text string, chatID int, username string) {
	// TODO make it look like normal
	msg := "Твої ігри:"
	var games []*storage.SantaUser
	storage.DB.Table("santa_users").Where("username = ?", username).Find(&games)
	MyGamesKeyboard := telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{},
			{*ButtonMain},
		},
		ResizeKeyboard: true,
	}
	for i := 0; i < len(games); i++ {
		if games[i].Game != "" {
			buttonName := fmt.Sprintf("Налаштування гри: %s (ID:%v)", games[i].Game, games[i].SantaID)
			button := NewButton(buttonName)
			AddButtonToKeyboard(button, &MyGamesKeyboard, i)
			msg = fmt.Sprintf("%s\n%s", msg, games[i].Game)
		}
	}
	p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msg, KeyboardReply: &MyGamesKeyboard})
	FSM.SetState(*MyGamesSate)
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
	for _, match := range storage.ListOfWishUpdates.Wishes {
		fmt.Printf("here --- %+v", match)
		if match.Username == username {
			match.Wish = text

			storage.DB.AddOrUpdateWishes(username)
			p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgWishesAdded})
		}
	}
}

func ExtractIDFromStringSettings(str string) string {
	// Налаштування гри: Ім‘я (ID:12345)
	var re = regexp.MustCompile(`(?m)ID:[0-9]+\)`)
	var id string
	for _, match := range re.FindAllString(str, -1) {
		id = strings.Split(match, ":")[1]
		id = strings.ReplaceAll(id, ")", "")
	}
	return id
}

func (p *Processor) ChooseTheGame(text string, chatID int, username string) {
	if len(text) > 17 {
		asRunes := []rune(text)
		reqStr := string(asRunes[:17])
		if reqStr == "Налаштування гри:" {
			fmt.Println("Yes")
		}
	}
	id := ExtractIDFromStringSettings(text)
	var games []*storage.SantaUser
	storage.DB.Table("santa_users").Where("username = ?", username).Find(&games)
	msg := fmt.Sprintf("Налаштування гри '%s'", games[0].Game)
	addWishesButton := telegram.InlineKeyboardButton{
		Text:         cmdChangeWishes,
		CallbackData: "change_wishes " + id,
	}
	keyboard := &telegram.InlineKeyboardMarkup{
		Buttons: [][]telegram.InlineKeyboardButton{{addWishesButton}},
	}
	p.tg.SendMessage(telegram.MessageParams{
		ChatID:         chatID,
		Text:           msg,
		KeyboardInline: keyboard,
	})
}

func cutTextToData(text string) (string, int) {
	i := strings.Index(text, " ")
	command := text[:i]
	idStr := text[i+1:]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic("error converting id to int")
	}
	return command, id
}
