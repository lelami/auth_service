package domain

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type role string

func RoleFromString(s string) (role, error) {
	switch role(s) {
	case UserRoleDefault:
		return UserRoleDefault, nil
	case UserRoleAdmin:
		return UserRoleAdmin, nil
	default:
		return "", errors.New("unknown role")
	}
}

const (
	UserRoleDefault role = "user"
	UserRoleAdmin        = "admin"
)

type User struct {
	ID       primitive.ObjectID `json:"id"`
	Login    string             `json:"login"`
	Password string             `json:"password"`
	Name     string             `json:"name"`
	Role     role               `json:"role"`
	Blocked  bool               `json:"blocked"`
}

type UserInfo struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}

type UserRole struct {
	ID   primitive.ObjectID `json:"id"`
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
