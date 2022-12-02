package telegram

import (
	"errors"
	"fmt"
	"main/clients/telegram"
	"main/events"
)

type Processor struct {
	tg     *telegram.Client
	offset int
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType = errors.New("unknown meta type")
)

func New(client *telegram.Client) *Processor {
	return &Processor{
		tg: client,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get events %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.CallbackQuery:
		return p.processCallbackQuery(event)
	default:
		return fmt.Errorf("can't process message %w", ErrUknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message %w", err)
	}
	if err := p.doMessage(event.Text, meta.ChatID, meta.Username); err != nil {
		return fmt.Errorf("can't process message %w", err)
	}
	return nil
}

func (p *Processor) processCallbackQuery(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message %w", err)
	}
	if err := p.doCallbackQuerry(event.Text, meta.ChatID, meta.Username); err != nil {
		return fmt.Errorf("can't process message %w", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("can't get meta %w", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	if updType == events.CallbackQuery {
		res.Meta = Meta{
			ChatID:   upd.CallbackQuerry.Message.Chat.ID,
			Username: upd.CallbackQuerry.From.Username,
		}
	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message != nil {
		return upd.Message.Text
	}
	if upd.CallbackQuerry != nil {
		return upd.CallbackQuerry.Data
	}
	return ""
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	}
	if upd.CallbackQuerry != nil {
		return events.CallbackQuery
	}
	return events.Unknown

}
