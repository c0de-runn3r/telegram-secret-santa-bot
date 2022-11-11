package telegram

import (
	"log"
	"strings"

	. "main/fsm"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%s'", text, username)
	state := FSM.CurrentState

	switch state {
	case *HelloState:
		return p.tg.SendMessage(chatID, "Hello")
	case *ByeState:
		return p.tg.SendMessage(chatID, "Buy")
	default:
		FSM.SetState(*HelloState)
		return p.tg.SendMessage(chatID, "No state set")
	}
}
