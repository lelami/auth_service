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

	userID, err := p.storage.GetUserByTgLink(username)
	if err != nil {
		log.Printf("user with tg link '%s' not found: %v", username, err)
		return err
	}

	user, err := p.storage.GetUser(*userID)
	if err != nil {
		log.Printf("user with id '%s' not found: %v", userID.Hex(), err)
		return err
	}

	user.ChatID = chatIDStr

	p.storage.SetUser(user)
	p.tg.SendMessage(chatID, msgHello)
	p.tg.SendMessage(chatID, msgSaved)
	return nil
}
