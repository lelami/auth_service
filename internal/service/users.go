package service

import (
	"authservice/internal/domain"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(lp *domain.LoginPassword) (*domain.UserToken, error) {

	if _, ok := users.CheckExistLogin(lp.Login); ok {
		return nil, errors.New("login " + lp.Login + " already exists")
	}

	newUser := domain.User{
		ID:       primitive.NewObjectID(),
		Login:    lp.Login,
		Password: hash(lp.Password),
		Role:     domain.UserRoleDefault,
	}

	if err := users.SetUser(&newUser); err != nil {
		return nil, err
	}

	token := createToken(lp.Login)

	if err := tokens.SetUserToken(token, newUser.ID); err != nil {
		return nil, err
	}

	return &domain.UserToken{
		UserId: newUser.ID,
		Token:  token,
	}, nil
}

func SignIn(lp *domain.LoginPassword) (*domain.UserToken, error) {

	userId, ok := users.CheckExistLogin(lp.Login)
	if !ok {
		return nil, errors.New("user not found")
	}

	user, err := users.GetUser(*userId)
	if err != nil {
		return nil, err
	}

	if user.Password != hash(lp.Password) {
		return nil, errors.New("wrong password")
	}

	token := createToken(lp.Login)

	if err := tokens.SetUserToken(token, *userId); err != nil {
		return nil, err
	}

	return &domain.UserToken{
		UserId: *userId,
		Token:  token,
	}, nil
}

func SetUserInfo(ui *domain.UserInfo) error {

	user, err := users.GetUser(ui.ID)
	if err != nil {
		return err
	}

	user.Name = ui.Name

	return users.SetUser(user)

}

func ChangePsw(up *domain.UserPassword) error {

	user, err := users.GetUser(up.ID)
	if err != nil {
		return err
	}

	user.Password = hash(up.Password)

	return users.SetUser(user)
}

func GetUserShortInfo(id primitive.ObjectID) (*domain.UserInfo, error) {

	user, err := users.GetUser(id)
	if err != nil {
		return nil, err
	}

	ui := domain.UserInfo{
		ID:              user.ID,
		Name:            user.Name,
		TelegramAccount: user.TelegramAccount,
	}

	return &ui, nil
}

func GetUserFullInfo(id primitive.ObjectID) (*domain.User, error) {

	user, err := users.GetUser(id)
	return user, err
}

func GetUserIDByToken(token string) (*primitive.ObjectID, error) {
	return tokens.GetUserByToken(token)
}

func hash(str string) string {
	hp := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hp[:])
}

func createToken(login string) string {

	timeChs := md5.Sum([]byte(time.Now().String()))
	loginChs := md5.Sum([]byte(login))

	return hex.EncodeToString(timeChs[:]) + hex.EncodeToString(loginChs[:])
}

func GetTelegramConnectLink(id primitive.ObjectID) (*domain.TelegramConnectLink, error) {

	_, err := users.GetUser(id)
	if err != nil {
		return nil, err
	}

	link := telegramBotService.BotUrl + "?start=" + id.Hex()

	return &domain.TelegramConnectLink{Link: link}, nil
}

func ConnectTelegramAccount(userId primitive.ObjectID, account *domain.TelegramAccount) error {

	user, err := users.GetUser(userId)
	if err != nil {
		return err
	}

	user.TelegramAccount = account

	return nil
}

func SendTelegramAuthCode(lp *domain.SendTelgramAuthCode) error {

	userId, ok := users.CheckExistLogin(lp.Login)
	if !ok {
		return errors.New("user not found")
	}

	user, err := users.GetUser(*userId)
	if err != nil {
		return err
	}

	if !user.IsTelegramConnected() {
		return errors.New("telegram is not connected")
	}

	code := rand.IntN(899999) + 100000 // radmon int from 100000 to 999999
	err = telegramAuthCodes.SetUserTelegramAuthCode(code, *userId)
	if err != nil {
		return err
	}

	go func(code int) {
		text := fmt.Sprintf("Ваш код авторизации: %d", code)
		telegramBotService.SendMessage(user.TelegramAccount, text)
	}(code)

	return nil
}

func SignInByTelegram(lp *domain.LoginTelegram) (*domain.UserToken, error) {

	userId, ok := users.CheckExistLogin(lp.Login)
	if !ok {
		return nil, errors.New("user not found")
	}

	user, err := users.GetUser(*userId)
	if err != nil {
		return nil, err
	}

	if !user.IsTelegramConnected() {
		return nil, errors.New("telegram is not connected")
	}

	authCode, err := telegramAuthCodes.GetTelegramAuthCodeByUserId(*userId)
	if err != nil {
		return nil, err
	}

	if authCode != lp.Code {
		return nil, errors.New("wrong code")
	}

	token := createToken(lp.Login)

	if err := tokens.SetUserToken(token, *userId); err != nil {
		return nil, err
	}

	if err := telegramAuthCodes.DeleteUserTelegramAuthCode(*userId); err != nil {
		return nil, err
	}

	return &domain.UserToken{
		UserId: *userId,
		Token:  token,
	}, nil
}
