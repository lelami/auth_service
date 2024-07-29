package httphandler

import (
	"authservice/internal/config"
	"authservice/internal/domain"
	"authservice/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input domain.LoginPassword
	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	userToken, err := service.SignUp(&input)
	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		respBody.SetError(err)
		return
	}

	respBody.SetData(userToken)
}

func SignIn(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input domain.LoginPassword
	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	userToken, err := service.SignIn(&input)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
		return
	}

	respBody.SetData(userToken)
}

func GetUserInfo(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	userID, _ := primitive.ObjectIDFromHex(req.Header.Get(HeaderUserID))

	info, err := service.GetUserShortInfo(userID)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
	}

	respBody.SetData(info)
}

func SetUserInfo(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input SetUserInfoReq

	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	userID, _ := primitive.ObjectIDFromHex(req.Header.Get(HeaderUserID))

	if err := service.SetUserInfo(&domain.UserInfo{
		ID:   userID,
		Name: input.Name,
	}); err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
		return
	}
}
func ChangePsw(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input ChangePswReq

	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	userID, _ := primitive.ObjectIDFromHex(req.Header.Get(HeaderUserID))
	err := service.ChangePsw(&domain.UserPassword{
		ID:       userID,
		Password: input.Password,
	})
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
		return
	}
}

// Of course, it would be logical to expand SetUserInfo,
// but we stick to business tasks.
func SetUserTgLink(resp http.ResponseWriter, req *http.Request) {

	cfg := config.GetConfig()

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input SetUserTgLinkReq

	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	normalizedUsername := normalizeTgLink(input.TgLink)
	fmt.Printf("Normalized username: %s\n", normalizedUsername)
	if normalizedUsername == "" {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
	}

	userID, _ := primitive.ObjectIDFromHex(req.Header.Get(HeaderUserID))
	if err := service.SetUserTgLink(&domain.UserTgLink{
		UserID: userID,
		TgLink: normalizedUsername,
	}); err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
		return
	}

	// We redirect to the bot so that the user starts chatting with it
	// This is due to the safety of TG
	http.Redirect(resp, req, cfg.BotLink, http.StatusSeeOther)
}

func ResetPsw(resp http.ResponseWriter, req *http.Request) {
	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	userID, _ := primitive.ObjectIDFromHex(req.Header.Get(HeaderUserID))
	err := service.ResetPsw(userID)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
		return
	}
}

func TgSignIn(resp http.ResponseWriter, req *http.Request) {
	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input domain.Login
	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	service.TgSignIn(&input)
}

func TgCheckOTP(resp http.ResponseWriter, req *http.Request) {
	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input domain.Code
	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	userID, _ := primitive.ObjectIDFromHex(req.Header.Get(HeaderUserID))
	userToken, err := service.TgCheckOTP(userID, &input)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
		return
	}

	respBody.SetData(userToken)
}

func readBody(req *http.Request, s any) error {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, s)
}
