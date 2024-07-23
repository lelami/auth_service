package httphandler

import (
	"authservice/internal/domain"
	"authservice/internal/service"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AdminGetUserInfo(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	id := req.URL.Query().Get("user_id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	info, err := service.GetUserFullInfo(userID)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
	}

	respBody.SetData(info)
}

func AdminBlockUser(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input BlockUserReq

	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	err := service.BlockUserByAdmin(&domain.UserBlocked{
		ID:      input.ID,
		Blocked: input.Blocked,
	})
	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		respBody.SetError(err)
		return
	}
}
