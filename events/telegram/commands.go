package telegram

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	case "all_players":
		FSM.SetState(*ActionState)
		p.AllPlayers(id, chatID, username)
	case "start_game":
		FSM.SetState(*ActionState)
		p.StartGame(id, chatID, username)
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
	id := storage.DB.AddNewGame(gameName, username, chatID)
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
		storage.DB.AddUserToGame(&game, username, chatID)
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
		Keyboard:        make([][]telegram.KeyboardButton, len(games)+1),
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	for i := 0; i < len(games); i++ {
		if games[i].Game != "" {
			buttonName := fmt.Sprintf("Налаштування гри: %s (ID:%v)", games[i].Game, games[i].SantaID)
			button := NewButton(buttonName)
			AddButtonToKeyboard(button, &MyGamesKeyboard, i)
			msg = fmt.Sprintf("%s\n%s", msg, games[i].Game)
		}
	}
	AddButtonToKeyboard(ButtonMain, &MyGamesKeyboard, len(games))
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
	id := ExtractIDFromStringSettings(text) // TODO probably need to change id extraction not from last command, but from message to what button connected
	var games []*storage.SantaUser
	storage.DB.Table("santa_users").Where("username = ?", username).Find(&games)
	msg := fmt.Sprintf("Налаштування гри '%s'", games[0].Game)
	showAllPlayersButton := telegram.InlineKeyboardButton{
		Text:         cmdShowAllPlayers,
		CallbackData: "all_players " + id,
	}
	addWishesButton := telegram.InlineKeyboardButton{
		Text:         cmdChangeWishes,
		CallbackData: "change_wishes " + id,
	}
	startGameButton := telegram.InlineKeyboardButton{
		Text:         cmdStartGame,
		CallbackData: "start_game " + id,
	}

	keyboard := &telegram.InlineKeyboardMarkup{
		Buttons: [][]telegram.InlineKeyboardButton{{showAllPlayersButton}, {addWishesButton}, {startGameButton}},
	}
	p.tg.SendMessage(telegram.MessageParams{
		ChatID:         chatID,
		Text:           msg,
		KeyboardInline: keyboard,
	})
}
func (p *Processor) AllPlayers(gameID int, chatID int, username string) {
	users, err := storage.DB.QueryAllPlayers(gameID)
	if err != nil {
		panic("no users found in this game")
	}
	resp := fmt.Sprintln("Список учасників:")
	for _, user := range users {
		resp = fmt.Sprintf("%s@%s\n", resp, user.Username)
	}
	p.tg.SendMessage(telegram.MessageParams{
		ChatID:        chatID,
		Text:          resp,
		KeyboardReply: &ActionKeyboard,
	})
}

func (p *Processor) StartGame(gameID int, chatID int, username string) { // TODO this piece of code is a piece of sh*t. need to change it
	admin, _ := storage.DB.QueryAdmin(gameID)
	if username != admin {
		msgIsNotAdmin := fmt.Sprintf("У вас немає доступу до цієї команди.\nПочати гру може лише @%s", admin)
		p.tg.SendMessage(telegram.MessageParams{
			ChatID: chatID,
			Text:   msgIsNotAdmin,
		})
		return
	}
	list, _ := storage.DB.QueryAllPlayers(gameID)
	if len(list) < 3 {
		p.tg.SendMessage(telegram.MessageParams{
			ChatID: chatID,
			Text:   "Кількість учасників має бути не менше 3",
		})
		return
	}
	res := DistributeSantas(gameID)
	for k, v := range res {
		msg := fmt.Sprintf("Ти даруєш подарунок @%s", v.Username)
		p.tg.SendMessage(telegram.MessageParams{
			ChatID: k.ChatID,
			Text:   msg,
		})
	}
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

func DistributeSantas(gameID int) map[storage.SantaUser]storage.SantaUser {
	list, _ := storage.DB.QueryAllPlayers(gameID)
	players := make([]storage.SantaUser, len(list))
	users := make([]storage.SantaUser, len(list))
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })
	copy(players, list)
	copy(users, list)
	user_pairs := make(map[storage.SantaUser]storage.SantaUser, len(list))
	for {
		if len(players) > 1 {
			if players[0] != users[0] {
				user_pairs[players[0]] = users[0]
				players = players[1:]
				users = users[1:]

			}
			if players[0] == users[0] {
				rand.Seed(time.Now().UnixNano())
				rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })
			}
		}
		if len(players) == 1 {
			user_pairs[players[0]] = users[0]
			return user_pairs
		}
	}
}
