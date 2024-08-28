package telegram

import (
	"log"
	"strconv"
	"strings"
)

const (
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	switch text {
	case StartCmd:
		return p.saveChatID(username, chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) saveChatID(username string, chatID int) error {
	chatIDStr := strconv.Itoa(chatID)

	if err := p.storage.SetUserChatID(username, chatIDStr); err != nil {
		p.tg.SendMessage(chatID, msgUserNotFound)
		return err
	}

	p.tg.SendMessage(chatID, msgHello)
	p.tg.SendMessage(chatID, msgSaved)
	return nil
}
