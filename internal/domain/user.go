package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	UserRoleDefault = "user"
	UserRoleAdmin   = "admin"
)

type Edu struct {
	Type    string `json:"type" bson:"type"`
	DocType string `json:"doc_type" bson:"doc_type"`
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Login     string             `json:"login" bson:"login"`
	Password  string             `json:"password" bson:"password"`
	Name      string             `json:"name" bson:"name"`
	Role      string             `json:"role" bson:"role"`
	Skills    []string           `json:"skills" bson:"skills"`
	Education []Edu              `json:"education" bson:"education"`
	Created   time.Time          `json:"created" bson:"created"`
	Updated   time.Time          `json:"updated" bson:"updated"`
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
