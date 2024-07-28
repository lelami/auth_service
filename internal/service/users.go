package service

import (
	"authservice/internal/domain"
	"authservice/internal/repository/tokendb"
	"authservice/internal/repository/userdb"
	"authservice/internal/service/rand"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const defaultPasswordLength = 8

type PasswordSender interface {
	SendPassword(psw string) error
}

var users userdb.DB
var tokens tokendb.DB
var pswSender PasswordSender

func Init(userDB userdb.DB, tokenDB tokendb.DB, pswSnd PasswordSender) {
	users = userDB
	tokens = tokenDB
	pswSender = pswSnd
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
		ID:   user.ID,
		Name: user.Name,
	}

	return &ui, nil
}

func GetUserFullInfo(id primitive.ObjectID) (*domain.User, error) {

	user, err := users.GetUser(id)
	return user, err
}

func SetUserRole(id primitive.ObjectID, roleString string) error {
	role, err := domain.RoleFromString(roleString)
	if err != nil {
		return errors.New("invalid role")
	}

	u, err := users.GetUser(id)
	if err != nil {
		return err
	}

	u.Role = role
	return users.SetUser(u)
}

func BlockUser(id primitive.ObjectID) error {
	u, err := users.GetUser(id)
	if err != nil {
		return err
	}

	u.Blocked = true
	return users.SetUser(u)
}

func UnblockUser(id primitive.ObjectID) error {
	u, err := users.GetUser(id)
	if err != nil {
		return err
	}

	u.Blocked = false
	return users.SetUser(u)
}

func ResetPasswordUser(id primitive.ObjectID) error {
	u, err := users.GetUser(id)
	if err != nil {
		return err
	}

	psw := rand.String(defaultPasswordLength)
	u.Password = hash(psw)

	if err := pswSender.SendPassword(psw); err != nil {
		return err
	}

	return users.SetUser(u)
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
