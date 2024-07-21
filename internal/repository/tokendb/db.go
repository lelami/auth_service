package tokendb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DB interface {
	GetUserByToken(token string) (*primitive.ObjectID, error)
	SetUserToken(token string, userID primitive.ObjectID) error
}
