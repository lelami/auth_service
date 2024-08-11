package telegramauthcodedb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DB interface {
	GetTelegramAuthCodeByUserId(userID primitive.ObjectID) (int, error)
	SetUserTelegramAuthCode(code int, userID primitive.ObjectID) error
	DeleteUserTelegramAuthCode(userID primitive.ObjectID) error
}
