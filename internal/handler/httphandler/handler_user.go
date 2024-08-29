package httphandler

import (
	"authservice/internal/domain"
	"authservice/internal/service"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"io"
	"net/http"
)

func SignUp(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	allReqCount.Add(req.Context(), 1)
	signUpReqCount.Add(req.Context(), 1)

	// Создание контекста и span для всего обработчика
	_, span := tracer.Start(req.Context(), "SignUp")
	defer span.End()

	var input domain.LoginPassword
	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", err.Error()))
		return
	}

	span.AddEvent("validation_started")
	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		err := errors.New("invalid input")
		respBody.SetError(err)

		span.SetStatus(codes.Error, "Validation failed")
		span.RecordError(err)
		return
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("login", input.Login))

	userToken, err := service.SignUp(&input)
	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		respBody.SetError(err)

		span.SetStatus(codes.Error, "Sign up failed")
		span.RecordError(err)
		return
	}

	respBody.SetData(userToken)
	span.SetStatus(codes.Ok, "")
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

	allReqCount.Add(req.Context(), 1)

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

func readBody(req *http.Request, s any) error {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, s)
}
