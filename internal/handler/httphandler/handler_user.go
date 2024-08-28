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

// SignUp godoc
// @Tags Auth
// @Summary Регистрация пользователя
// @Param request body domain.LoginPassword true "логин и пароль нового пользователя"
// @Success 200 {object} domain.UserToken
// @Failure 400
// @Failure 404
// @Failure 409 {object} ErrorResponse
// @Failure 422
// @Failure 500
// @Router /sign_up [post]
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

// SignIn godoc
// @Tags Auth
// @Summary Вход пользователя
// @Param request body domain.LoginPassword true "логин и пароль нового пользователя"
// @Success 200 {object} domain.UserToken
// @Failure 400
// @Failure 404
// @Failure 422
// @Failure 500
// @Router /sign_in [post]
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

// GetUserInfo godoc
// @Tags Auth
// @Summary Получение личных данных пользователя
// @Param Authorization header string true "токен пользователя"
// @Success 200 {object} domain.UserInfo
// @Failure 404
// @Failure 500
// @Router /get_user_info [get]
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

// SetUserInfo godoc
// @Tags Auth
// @Summary Изменение личных данных пользователя
// @Param Authorization header string true "токен пользователя"
// @Param request body SetUserInfoReq true "новое имя пользователя"
// @Success 200 {object} domain.UserInfo
// @Failure 400
// @Failure 404
// @Failure 422
// @Failure 500
// @Router /set_user_info [post]
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

// ChangePsw godoc
// @Tags Auth
// @Summary Изменение пороля пользователем
// @Param Authorization header string true "токен пользователя"
// @Param request body ChangePswReq true "новый пароль пользователя"
// @Success 200 {object} domain.UserPassword
// @Failure 400
// @Failure 404
// @Failure 422
// @Failure 500
// @Router /change_psw [post]
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

// SetUserTgLink godoc
// @Tags Auth
// @Summary Привязка профиля телеграм
// @Param Authorization header string true "токен пользователя"
// @Param request body SetUserTgLinkReq true "ссылка на телеграм пользователя"
// @Success 200 {object} RedirectTgLink
// @Failure 400
// @Failure 404
// @Failure 422
// @Failure 500
// @Router /set_tg_link [post]
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
	resp.WriteHeader(http.StatusOK)
	respBody.Data = RedirectTgLink{
		RedirectURL: cfg.BotLink,
	}
}

// ResetPsw godoc
// @Tags Auth
// @Summary Сброс пароля пользователем
// @Param Authorization header string true "токен пользователя"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 422
// @Failure 500
// @Router /reset_psw [post]
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

// TgSignIn godoc
// @Tags Auth
// @Summary Отправка проверочного кода для входа через ТГ
// @Param request body domain.Login true "логин пользователя"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 422
// @Failure 500
// @Router /sign_in_via_tg [post]
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

// TgCheckOTP godoc
// @Tags Auth
// @Summary Вход по коду через телеграм
// @Param request body domain.Code true "4ёх значный проверочный код"
// @Param User-ID header string true "ID пользователя"
// @Success 200 {object} domain.UserToken
// @Failure 400
// @Failure 404
// @Failure 422
// @Failure 500
// @Router /check_otp [post]
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

// GetUserWithRole godoc
// @Tags Auth
// @Summary Используется для получения информации о пользователе другими сервисами
// @Param Authorization header string true "токен пользователя"
// @Param X-Service-Key header string true "сервисный ключ"
// @Success 200 {object} domain.UserInfoWithRole
// @Failure 400
// @Failure 404
// @Failure 422
// @Failure 500
// @Router /get_service_user_info [get]
func GetUserWithRole(resp http.ResponseWriter, req *http.Request) {
	cfg := config.GetConfig()

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	userID, _ := primitive.ObjectIDFromHex(req.Header.Get(HeaderUserID))
	serviceKey := req.Header.Get(HeaderServiceKey)

	if serviceKey != cfg.ServiceKey || serviceKey == "" {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid service"))
		return
	}

	info, err := service.GetUserInfoWithRole(userID)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
	}

	respBody.SetData(info)
}

func readBody(req *http.Request, s any) error {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, s)
}
