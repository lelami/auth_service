package telegram

import (
	"authservice/internal/domain"
	e "authservice/internal/helpers"
	"authservice/internal/service"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

const (
	codeLength = 6

	StartCmd  = "/start"
	HelpCmd   = "/help"
	RepeatCmd = "/repeat"
)

func (fp *FetchProcessor) doCmd(cmd string, chatID int, userName string) error {
	cmd = strings.TrimSpace(cmd)

	log.Printf("got new cmd=%s chatID=%d from userName=%s", cmd, chatID, userName)

	//start:
	// help:
	switch cmd {
	case StartCmd:
		return fp.SendHello(chatID, userName)
	case HelpCmd:
		return fp.SendHelp(chatID)
	case RepeatCmd:
		return fp.SendRepeat(chatID, userName)
	default:
		return fp.tg.SendMessage(chatID, msgUnknownCmd)

	}
} /**/
func (fp *FetchProcessor) SendHelp(chatID int) error {
	return fp.tg.SendMessage(chatID, msgHelp)
}
func (fp *FetchProcessor) SendRepeat(chatID int, userName string) error {
	code, err := oneTimeCode()
	if err != nil {
		return e.WrapIfErr("can't send code", err)
	}
	err = service.SetCode(&domain.TgNameCode{
		Name: userName,
		Code: code,
	})
	switch {
	case errors.Is(err, os.ErrNotExist):
		return fp.tg.SendMessage(chatID, fmt.Sprintf(msgNotFount, userName))
	case err != nil:
		return e.WrapIfErr("can't send code", err)
	}
	return fp.tg.SendMessage(chatID, fmt.Sprintf("%s%s", msgRepeat, code))
}
func (fp *FetchProcessor) SendHello(chatID int, userName string) error {

	code, err := oneTimeCode()
	if err != nil {
		return e.WrapIfErr("can't send code", err)
	}
	err = service.SetCode(&domain.TgNameCode{
		Name: userName,
		Code: code,
	})
	switch {
	case errors.Is(err, os.ErrNotExist):
		return fp.tg.SendMessage(chatID, fmt.Sprintf(msgNotFount, userName))
	case err != nil:
		return e.WrapIfErr("can't send code", err)
	}

	return fp.tg.SendMessage(chatID, fmt.Sprintf("%s%s", msgHello, code))
}
func oneTimeCode() (code string, err error) {
	numberSet := "0123456789"
	var bsCode strings.Builder
	for i := 0; i < codeLength; i++ {
		random := rand.Intn(len(numberSet))
		_, err = bsCode.WriteString(string(numberSet[random]))
		if err != nil {
			return "", e.WrapIfErr("can't send one time code", err)
		}
	}
	return bsCode.String(), nil
}
