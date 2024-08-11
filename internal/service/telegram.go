package service

import (
	"authservice/internal/config"
	"authservice/internal/domain"
	"context"
	"log"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type TelegramService struct {
	Bot    *tgbotapi.BotAPI
	BotUrl string
}

func NewTelegramService(cfg *config.TelegramConfig) (*TelegramService, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = true // Включить/отключить дебаг
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &TelegramService{Bot: bot, BotUrl: cfg.BotUrl}, nil
}

func (s *TelegramService) SendMessage(account *domain.TelegramAccount, text string) error {
	msg := tgbotapi.NewMessage(account.ID, text)
	_, err := s.Bot.Send(msg)
	return err
}

func (s *TelegramService) ListenForUpdates(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.Bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				s.handleMessage(update.Message)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *TelegramService) handleMessage(message *tgbotapi.Message) {
	if strings.HasPrefix(message.Text, "/start") {
		re := regexp.MustCompile(`^/start (\S+)$`)
		matches := re.FindStringSubmatch(message.Text)
		if len(matches) > 1 {
			id := matches[1]
			s.handleStartCommand(message, id)
		}
	}
}

func (s *TelegramService) handleStartCommand(message *tgbotapi.Message, id string) {

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.sendReply(message.Chat.ID, "Ошибка подключения аккаунта, попробуйте снова")
		return
	}

	err = ConnectTelegramAccount(userId, getTelegramAccount(message))
	if err != nil {
		s.sendReply(message.Chat.ID, "Ошибка подключения аккаунта, попробуйте снова")
		return
	}

	s.sendReply(message.Chat.ID, "Аккаунт успешно подключён")
}

func (s *TelegramService) sendReply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := s.Bot.Send(msg)
	if err != nil {
		log.Printf("ERROR failed to send message: %v", err)
	}
}

func getTelegramAccount(message *tgbotapi.Message) *domain.TelegramAccount {
	return &domain.TelegramAccount{
		ID:        message.From.ID,
		FirstName: message.From.FirstName,
		LastName:  message.From.LastName,
	}
}
