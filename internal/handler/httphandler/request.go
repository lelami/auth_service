package httphandler

import "go.mongodb.org/mongo-driver/bson/primitive"

type SetUserInfoReq struct {
	Name string `json:"name"`
}

type ChangePswReq struct {
	Password string `json:"password"`
}

type BlockUserReq struct {
	ID      primitive.ObjectID `json:"id"`
	Blocked bool               `json:"blocked"`
}

func (r SetUserInfoReq) IsValid() bool {
	return r.Name != ""
}

func (r ChangePswReq) IsValid() bool {
	return r.Password != ""
}
