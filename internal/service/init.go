package service

import (
	"authservice/internal/repository/telegramauthcodedb"
	"authservice/internal/repository/tokendb"
	"authservice/internal/repository/userdb"
)

var users userdb.DB
var tokens tokendb.DB
var telegramAuthCodes telegramauthcodedb.DB
var telegramBotService TelegramService

func Init(userDB userdb.DB, tokenDB tokendb.DB, telegramTokenDB telegramauthcodedb.DB, telegramService TelegramService) {
	users = userDB
	tokens = tokenDB
	telegramAuthCodes = telegramTokenDB
	telegramBotService = telegramService
}
