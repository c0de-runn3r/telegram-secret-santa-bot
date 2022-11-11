package telegram

import (
	"errors"
	"fmt"
	"log"
	"main/storage"
	"main/storage/files"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}
	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}
	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return fmt.Errorf("can't do command: save page %w", err)
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}
	if err := p.storage.Save(page); err != nil {
		return fmt.Errorf("can't do command: save page %w", err)
	}
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return fmt.Errorf("can't do command: save page %w", err)
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, files.ErrNoSavedPages) {
		return fmt.Errorf("cant do ocommand: can't send random %w", err)
	}
	if errors.Is(err, files.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return fmt.Errorf("cant do ocommand: can't send random %w", err)
	}
	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
