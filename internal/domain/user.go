package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	ID      primitive.ObjectID `json:"id"`
	Blocked bool               `json:"blocked"`
}

type UserRole struct {
	ID   primitive.ObjectID `json:"id"`
	Role string             `json:"role"`
}

type UserTgLink struct {
	ID     primitive.ObjectID `json:"id"`
	TgLink string             `json:"tg_link"`
}

type TgMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type UserChatID struct {
	ID     primitive.ObjectID `json:"id"`
	TgLink string             `json:"tg_link"`
	ChatID string             `json:"chat_id"`
}
