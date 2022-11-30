package telegram

// outgoing messages from bot
const (
	msgHelp            = "Створи нову гру, або ж приєднайся до вже існуючої."
	msgHello           = "Привіт!\nЯ Таємний Санта. Я допоможу тобі створити свято!\n\n" + msgHelp
	msgUnknownCommand  = "Невідома команда!"
	msgNameNewGame     = "Надішли мені назву нової гри"
	msgSendIDOfGame    = "Надішли мені ID гри до якої хочеш приєднатись"
	msgCancel          = "Відміна. Повертаємось на головну"
	msgSendIntNotStr   = "ID складається лише з цифр"
	msgUndefinedGameID = "Такого ID гри не існує. Спробуйте ще раз."
	msgAddWishes       = "Напиши сюди одним повідомленням твої побажання щодо подарунків. Якщо потрібно вказати адресу для відправки, також додай сюди)"
	msgWishesAdded     = "Твої побажання оновлено"
)

// incoming commands from user

// comands with slash
const (
	StartCmd = "/start"
	HelpCmd  = "/help"
)

// commands without slash
const (
	cmdCreateNewGame         = "Створити нову гру"
	cmdConnectToExistingGame = "Приєднатись до гри"
	cmdMain                  = "На головну"
	cmdChangeWishes          = "Додати/змінити побажання"
)
