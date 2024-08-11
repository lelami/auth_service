package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserRoleDefault = "user"
	UserRoleAdmin   = "admin"
)

type User struct {
	ID              primitive.ObjectID `json:"id"`
	Login           string             `json:"login"`
	Password        string             `json:"password"`
	Name            string             `json:"name"`
	Role            string             `json:"role"`
	TelegramAccount *TelegramAccount   `json:"telegram_account"`
}

type TelegramAccount struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (user *User) IsTelegramConnected() bool {
	return user.TelegramAccount != nil
}

type UserInfo struct {
	ID                  primitive.ObjectID `json:"id"`
	Name                string             `json:"name"`
	TelegramAccount     *TelegramAccount   `json:"telegram_account"`
}

type UserPassword struct {
	ID       primitive.ObjectID `json:"id"`
	Password string             `json:"password"`
}

type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserToken struct {
	UserId primitive.ObjectID `json:"id"`
	Token  string             `json:"token"`
}

type TelegramConnectLink struct {
	Link string `json:"link"`
}

type SendTelgramAuthCode struct {
	Login string `json:"login"`
}

type LoginTelegram struct {
	Login string `json:"login"`
	Code  int    `json:"code"`
}
