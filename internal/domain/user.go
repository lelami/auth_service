package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserRoleDefault = "user"
	UserRoleAdmin   = "admin"
)

type User struct {
	ID       primitive.ObjectID `json:"id"`
	Login    string             `json:"login"`
	Password string             `json:"password"`
	Name     string             `json:"name"`
	Role     string             `json:"role"`
	Blocked  bool               `json:"blocked"`
	TgLink   string             `json:"tg_link"`
	ChatID   string             `json:"chat_id"`
}

type UserInfo struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}

type UserInfoWithRole struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
	Role string             `json:"role"`
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

type UserBlocked struct {
	UserID  primitive.ObjectID `json:"id"`
	Blocked bool               `json:"blocked"`
}

type UserRole struct {
	UserID primitive.ObjectID `json:"id"`
	Role   string             `json:"role"`
}

type UserTgLink struct {
	UserID primitive.ObjectID `json:"id"`
	TgLink string             `json:"tg_link"`
}

type UserChatID struct {
	UserID primitive.ObjectID `json:"id"`
	TgLink string             `json:"tg_link"`
	ChatID string             `json:"chat_id"`
}

type UserOTP struct {
	UserID    primitive.ObjectID `json:"id"`
	Code      string             `json:"code"`
	CreatedAt time.Time          `json:"created_at"`
	Expiry    time.Time          `json:"expiry"`
	Used      bool               `json:"used"`
}

type Login struct {
	Login string `json:"login"`
}

type Code struct {
	Code string `json:"code"`
}
