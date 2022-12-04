package telegram

// outgoing messages from bot
const (
	msgHelp            = "Створи нову гру, або ж приєднайся до вже існуючої.\nПоки що бот має неповний функціонал, протягом наступних днів буде доповнюватись. Проте основні функції вже працюють.\nЩодо будь-яких помилок чи побажань пишіть @wanna_run_around_earth"
	msgHello           = "Хо-хо-хо!\nЯ - Таємний Санта 🎅\nЯ допоможу тобі створити свято разом з твоїми друзями!\n\n" + msgHelp
	msgSendStart       = "Натисни /start"
	msgUnknownCommand  = "Невідома команда!"
	msgNameNewGame     = "Хо-хо-хо!\nНадішли мені назву нової гри"
	msgSendIDOfGame    = "Хо-хо-хо!\nНадішли мені ID гри до якої хочеш приєднатись"
	msgCancel          = "Відміна. Повертаємось на головну"
	msgSendIntNotStr   = "ID складається лише з цифр"
	msgUndefinedGameID = "Такого ID гри не існує. Спробуйте ще раз."
	msgAddWishes       = "Хо-хо-хо!\nНапиши сюди одним повідомленням твої побажання щодо подарунків. Якщо потрібно вказати адресу для відправки, також додай сюди)"
	msgWishesAdded     = "Хо-хо-хо!\nТвої побажання оновлено"
	msgAlreadyInGame   = "Ти вже приєднався/приєдналась до цієї гри"
	msgSmthWrong       = "Щось пішло не так..."
	msgUserDeleted     = "Ти покинув цю гру"
	msgGameDeleted     = "Гра та всі учасники видалені"
)

// incoming commands from user

// comands with slash
const (
	StartCmd = "/start"
	HelpCmd  = "/help"
)

// commands without slash
const (
	cmdCreateNewGame         = "Створити нову гру 🎁"
	cmdConnectToExistingGame = "Приєднатись до гри 🚪"
	cmdCheckMyGames          = "Переглянути мої ігри"
	cmdMain                  = "На головну ↩️"
	cmdChangeWishes          = "Додати/змінити побажання"
	cmdShowAllPlayers        = "Показати всіх учасників"
	cmdStartGame             = "Розпочати гру"
	cmdQuitGame              = "Покинути гру"
)
