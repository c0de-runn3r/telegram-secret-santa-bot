package telegram

type UpdateResponce struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard        bool               `json:"resize_keyboard"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard"`
	InputFieldPlaceholder string             `json:"input_field_placeholder"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type InlineKeyboardMarkup struct {
	Buttons [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
	Url          string `json:"url"`
}

type User struct {
	ID       int    `json:"id"`
	IsBot    bool   `json:"is_bot"`
	Name     string `json:"first_name"`
	Username string `json:"username"`
}

type CallbackQuery struct {
	ID   string `json:"id"`
	From User   `json:"from"`
	Data string `json:"callback_data"`
}
