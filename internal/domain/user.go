package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	UserRoleDefault = "user"
	UserRoleAdmin   = "admin"
)

type User struct {
	ID       primitive.ObjectID `json:"id"`
	Login    string             `json:"login"`
	Password string             `json:"-"`
	Name     string             `json:"name"`
	Role     string             `json:"role"`
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
