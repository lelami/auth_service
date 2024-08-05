package service

import (
	"authservice/internal/domain"
	"authservice/internal/repository/tokendb"
	"authservice/internal/repository/userdb"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

// const botName = "https://t.me/OneTimeAuthBot"
var botName string
var users userdb.DB
var tokens tokendb.DB

func Init(userDB userdb.DB, tokenDB tokendb.DB, bName string) {
	users = userDB
	tokens = tokenDB
	botName = bName
}

func SignInByTelegram() (string, error) {
	return botName, nil
}
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
func SignInByCode(lp *domain.LoginCode) (*domain.UserToken, error) {
	userId, ok := users.CheckExistLogin(lp.Login)
	if !ok {
		return nil, errors.New("user not found")
	}

	user, err := users.GetUser(*userId)
	if err != nil {
		return nil, err
	}

	if user.OneTimeCode != hash(lp.Code) {
		return nil, errors.New("wrong code")
	}

	token := createToken(lp.Login)

	if err := tokens.SetUserToken(token, *userId); err != nil {
		return nil, err
	}
	user.OneTimeCode = ""

	if err = users.SetUser(user); err != nil {
		return nil, err
	}
	return &domain.UserToken{
		UserId: *userId,
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
func SetTelegramInfo(tgInfo *domain.TelegramInfo) error {
	user, err := users.GetUser(tgInfo.ID)
	if err != nil {
		return err
	}
	user.TelegramName = tgInfo.Name

	return users.SetUser(user)
}

func SetCode(up *domain.TgNameCode) error {

	userId, ok := users.CheckExistTgName(up.Name)
	if !ok {
		return os.ErrNotExist
	}
	user, err := users.GetUser(*userId)
	if err != nil {
		return err
	}
	user.OneTimeCode = hash(up.Code)

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
		ID:   user.ID,
		Name: user.Name,
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
