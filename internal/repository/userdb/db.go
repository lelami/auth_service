package userdb

import (
	"authservice/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DB interface {
	CheckExistLogin(login string) (*primitive.ObjectID, bool)
	GetUser(id primitive.ObjectID) (*domain.User, error)
	SetUser(user *domain.User) error
	SetUserTgLink(utg *domain.UserTgLink) error
	GetUserByTgLink(tgLink string) (*primitive.ObjectID, error)
	CheckExistChatID(id primitive.ObjectID) (*string, bool)
}
