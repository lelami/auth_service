package telegram

import (
	"authservice/internal/clients/events"
	"authservice/internal/clients/telegram"
	e "authservice/internal/helpers"
	"errors"
)

var (
	ErrUnknowEventType = errors.New("Unknow event type")
	ErrUnknowMetaType  = errors.New("Unknow meta type")
)

type FetchProcessor struct {
	tg     *telegram.Client
	offset int
}
type Meta struct {
	ChatID   int
	UserName string
}

func New(tg *telegram.Client) *FetchProcessor {
	return &FetchProcessor{
		tg: tg,
	}
}
func (fp *FetchProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := fp.tg.Updates(fp.offset, limit)
	if err != nil {
		return nil, e.WrapIfErr("can't get event", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}
	res := make([]events.Event, 0, len(updates))
	for _, update := range updates {
		res = append(res, event(update))
	}
	fp.offset = updates[len(updates)-1].ID + 1
	return res, nil
}
func (fp *FetchProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return fp.processMessage(event)
	default:
		return ErrUnknowEventType
	}
}
func (fp *FetchProcessor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.WrapIfErr("can't process mesage", err)
	}
	if err := fp.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return e.WrapIfErr("can't process message", err)
	}
	return nil
}
func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.WrapIfErr("can't get meta from event", ErrUnknowMetaType)
	}
	return res, nil
}
func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.UserName,
		}
	}
	return res
}
func fetchText(upd telegram.Update) string {
	if nil == upd.Message {
		return ""
	}
	return upd.Message.Text
}
func fetchType(u telegram.Update) events.Type {
	if nil == u.Message {
		return events.Unknow
	}
	return events.Message
}
