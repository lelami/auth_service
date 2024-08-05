package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	UserRoleDefault = "user"
	UserRoleAdmin   = "admin"
)

type User struct {
	ID           primitive.ObjectID `json:"id"`
	Login        string             `json:"login"`
	Password     string             `json:"password"`
	Name         string             `json:"name"`
	Role         string             `json:"role"`
	TelegramName string             `json:"tg_name"`
	OneTimeCode  string             `json:"code"`
}

type UserInfo struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}
type TelegramInfo struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
	Code string             `json:"code"`
}
type UserPassword struct {
	ID       primitive.ObjectID `json:"id"`
	Password string             `json:"password"`
}

type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type LoginCode struct {
	Login string `json:"login"`
	Code  string `json:"code"`
}
type TgNameCode struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type UserToken struct {
	UserId primitive.ObjectID `json:"id"`
	Token  string             `json:"token"`
}
