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
	userFSM := FindOrCreateUsersFSM(username)
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%s'", text, username)
	ok, startID := checkIfStartHasID(text)
	if ok {
		userFSM.SetState(*ActionState)
		p.ConnectToExistingGame(startID, chatID, username)
		return nil
	}
	state := userFSM.CurrentState
	keyboard := makeActionKeyboard(username)
	fmt.Printf("%+v", keyboard)
	switch text { //for commands
	case StartCmd: // /start
		userFSM.SetState(*ActionState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgHello})
	case HelpCmd: // /help
		userFSM.SetState(*ActionState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgHelp, KeyboardReply: makeActionKeyboard(username)})
	case cmdMain: // to the main menu
		userFSM.SetState(*ActionState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgCancel, KeyboardReply: makeActionKeyboard(username)})
	default: // other cases
		switch state {
		case *ActionState: // actions menu
			p.ProcessAction(text, chatID, username)
		case *NewGameNameState: // receive name of the new game
			userFSM.SetState(*ActionState)
			p.CreateNewGame(text, chatID, username)
		case *MyGamesSate: // receive id to change settings of the game
			p.ChooseTheGame(text, chatID, username)
		case *UpdateWishesState: // receive wishes text to update wishes
			p.UpdateWishes(text, chatID, username)
			userFSM.SetState(*ActionState)
		case *BudgetState: // receive budget sum
			p.ChangeBudget(text, chatID, username)
			userFSM.SetState(*ActionState)
		default:
			p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgUnknownCommand})
		}
	}

	return nil
}

func (p *Processor) doCallbackQuerry(text string, chatID int, username string) error {
	log.Printf("got new callback data '%s' from '%s'", text, username)
	userFSM := FindOrCreateUsersFSM(username)
	command, id := cutTextToData(text)
	switch command {
	case "change_wishes": // create struct with game id and nickname and go to state which receives wishes
		userFSM.SetState(*UpdateWishesState)
		storage.ListOfWishUpdates.Wishes = append(storage.ListOfWishUpdates.Wishes, &storage.WishUpdateInfo{
			ID:       id,
			Username: username})
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgAddWishes})
	case "all_players": // show list of players
		userFSM.SetState(*ActionState)
		p.AllPlayers(id, chatID, username)
	case "start_game": // roll the list
		userFSM.SetState(*ActionState)
		p.StartGame(id, chatID, username)
	case "quit_game": // leave the game
		userFSM.SetState(*ActionState)
		p.QuitGame(id, chatID, username)
	case "change_budget": // change budget
		userFSM.SetState(*BudgetState)
		storage.ListOfBudgetUpdates.Budgets = append(storage.ListOfBudgetUpdates.Budgets, &storage.BudgetInfo{
			ID:       id,
			Username: username,
		})
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgSendBudget})
	default:
		p.tg.SendMessage(telegram.MessageParams{
			ChatID: chatID,
			Text:   msgSmthWrong,
		})
	}

	return nil
}

func (p *Processor) ProcessAction(text string, chatID int, username string) {
	userFSM := FindOrCreateUsersFSM(username)
	if len(text) > 17 {
		asRunes := []rune(text)
		reqStr := string(asRunes[:17])
		if reqStr == "Налаштування гри:" {
			p.ChooseTheGame(text, chatID, username)
			return
		}
	}
	switch text {
	case cmdCreateNewGame:
		userFSM.SetState(*NewGameNameState)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgNameNewGame})
	default:
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgUnknownCommand})
	}

}

func (p *Processor) CreateNewGame(gameName string, chatID int, username string) {
	userFSM := FindOrCreateUsersFSM(username)
	log.Printf("creating new game [%s]", gameName)
	id := storage.DB.AddNewGame(gameName, username, chatID)
	msg := fmt.Sprintf("Хо-хо-хо!\nНову гру %s створено.\nID: %v\nТепер в налаштуваннях гри ти можеш змінити бюджет, побажання, і отримати доступ до адмінських функцій.\nПерешли наступне повідомлення своїм друзям, щоб вони могли приєднатись.", gameName, id)
	msg2 := fmt.Sprintf("Хо-хо-хо!\nЗапрошую тебе до гри в Таємного Санту🎅\nПереходь за цим посиланням:\nhttps://t.me/SecretSantaUkrBot?start=%v", id)
	p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msg})
	p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msg2, KeyboardReply: makeActionKeyboard(username)})
	userFSM.SetState(*ActionState)
}

// TODO budget feature

func (p *Processor) ConnectToExistingGame(strID string, chatID int, username string) {
	userFSM := FindOrCreateUsersFSM(username)
	gameID, err := strconv.Atoi(strID)
	if err != nil {
		log.Println("Can't convert stringID into int")
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgSendIntNotStr})
		userFSM.SetState(*ConnectExistingGameState)
		return
	}
	var game storage.Game
	storage.DB.First(&game, gameID)
	if game.ID != 0 {
		players, _ := storage.DB.QueryAllPlayers(gameID)
		for _, player := range players {
			if username == player.Username {
				p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgAlreadyInGame, KeyboardReply: makeActionKeyboard(username)})
				userFSM.SetState(*ActionState)
				return
			}
		}
		storage.DB.AddUserToGame(&game, username, chatID)
		msg := fmt.Sprintf("Хо-хо-хо!\nТи приєднався до %s\nЩасливого Різдва!\nНе забудь додати побажання 🎁\nБюджет цієї гри: %s\n Це можна зробити в налаштуваннях цієї гри ", game.Name, game.Budget)
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msg, KeyboardReply: makeActionKeyboard(username)})
		userFSM.SetState(*ActionState)
	} else {
		p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgUndefinedGameID})
		userFSM.SetState(*ConnectExistingGameState)
	}
}

func (p *Processor) ChangeGameSettings(text string, chatID int, username string) {
	userFSM := FindOrCreateUsersFSM(username)
	switch text {
	case cmdChangeWishes:
		userFSM.SetState(*UpdateWishesState)
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
			p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgWishesAdded, KeyboardReply: makeActionKeyboard(username)})
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

func checkIfStartHasID(text string) (bool, string) {
	if len(text) > 10 {
		if text[:6] == "/start" {
			idStr := text[7:]
			return true, idStr
		}
	}
	return false, ""
}

func (p *Processor) ChooseTheGame(text string, chatID int, username string) {
	if len(text) > 17 {
		asRunes := []rune(text)
		reqStr := string(asRunes[:17])
		if reqStr != "Налаштування гри:" {
			p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgUnknownCommand, KeyboardReply: makeActionKeyboard(username)})
			return
		}
	}
	id := ExtractIDFromStringSettings(text)
	keyboard := createSettingsKeyboard(text, username, id)
	var game *storage.SantaUser
	storage.DB.Table("santa_users").Where("santa_id = ?", id).First(&game)
	msg := fmt.Sprintf("Ельфи готові виконати будь-яку твою забаганку!\n⚙️ Налаштування гри %s", game.Game)

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
	resp := "📃 Список Сант, а також тих хто чекає своїх подаруночків:\n"

	for _, user := range users {
		resp = fmt.Sprintf("%s@%s\n", resp, user.Username)
	}

	p.tg.SendMessage(telegram.MessageParams{
		ChatID:        chatID,
		Text:          resp,
		KeyboardReply: makeActionKeyboard(username),
	})
}

func (p *Processor) StartGame(gameID int, chatID int, username string) {
	admin, _ := storage.DB.QueryAdmin(gameID)
	if username != admin {
		msgIsNotAdmin := fmt.Sprintf("Ельфи ще очікують списки подарунків!\nПочати гру може лише головний Санта @%s", admin)
		p.tg.SendMessage(telegram.MessageParams{
			ChatID: chatID,
			Text:   msgIsNotAdmin,
		})
		return
	}
	var game storage.Game
	storage.DB.Table("games").Where("id = ?", gameID).First(&game)
	if game.Rolled {
		p.tg.SendMessage(telegram.MessageParams{
			ChatID: chatID,
			Text:   msgGameIsBeenRolled,
		})
		return
	}
	list, _ := storage.DB.QueryAllPlayers(gameID)
	if len(list) < 3 {
		p.tg.SendMessage(telegram.MessageParams{
			ChatID: chatID,
			Text:   "Кількість Сант має бути не менше 3",
		})
		return
	}
	res := DistributeSantas(gameID)
	for k, v := range res {
		msg := fmt.Sprintf("Хо-хо-хо! Різдвяне чудо!❄️\nГра %s розпочалась!\nТепер ти - Санта🎅 для @%s\nЙого побажання🎁 такі:\n%s\nНагадаю, що бюджет гри: %s", game.Name, v.Username, v.Wishes, game.Budget)
		p.tg.SendMessage(telegram.MessageParams{
			ChatID: k.ChatID,
			Text:   msg,
		})
	}
}

func (p *Processor) QuitGame(gameID int, chatID int, username string) {
	admin, _ := storage.DB.QueryAdmin(gameID)
	if username == admin {
		storage.DB.DeleteGameAndAllUsers(gameID)
		p.tg.SendMessage(telegram.MessageParams{
			ChatID:        chatID,
			Text:          msgGameDeleted,
			KeyboardReply: makeActionKeyboard(username),
		})
	} else {
		storage.DB.DeleteUserFromGame(username, gameID)
		p.tg.SendMessage(telegram.MessageParams{
			ChatID:        chatID,
			Text:          msgUserDeleted,
			KeyboardReply: makeActionKeyboard(username),
		})
	}
}
func (p *Processor) ChangeBudget(budget string, chatID int, username string) {
	for _, match := range storage.ListOfBudgetUpdates.Budgets {
		if match.Username == username {
			match.Budget = budget
			storage.DB.AddOrUpdateBudget(username)
			p.tg.SendMessage(telegram.MessageParams{ChatID: chatID, Text: msgBudgetUpdated, KeyboardReply: makeActionKeyboard(username)})
		}
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
	storage.DB.Table("games").Where("id = ?", gameID).Update("rolled", true)
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
