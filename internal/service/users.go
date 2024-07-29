package service

import (
	"authservice/internal/domain"
	"authservice/internal/repository/tokendb"
	"authservice/internal/repository/userdb"
	"authservice/internal/server"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var users userdb.DB
var tokens tokendb.DB

func Init(userDB userdb.DB, tokenDB tokendb.DB) {
	users = userDB
	tokens = tokenDB
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
		Blocked:  false,
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

func BlockUserByAdmin(ub *domain.UserBlocked) error {

	user, err := users.GetUser(ub.ID)
	if err != nil {
		return err
	}

	user.Blocked = ub.Blocked

	return users.SetUser(user)
}

func ChangeRole(ur *domain.UserRole) error {

	user, err := users.GetUser(ur.ID)
	if err != nil {
		return err
	}

	user.Role = ur.Role

	return users.SetUser(user)
}

func ResetPsw(id primitive.ObjectID) error {

	user, err := users.GetUser(id)
	if err != nil {
		return err
	}

	newPsw, err := randPsw()
	if err != nil {
		fmt.Println("Error generating password:", err)
	}

	user.Password = hash(newPsw)

	// печатаем в лог
	log.Printf("Updated psw %s -> ID %s", newPsw, id.Hex())

	// или же отправляем в ТГ
	chatID, ok := users.CheckExistChatID(id)
	if ok && chatID != nil {
		chatIDInt, err := strconv.Atoi(*chatID)
		if err != nil {
			log.Println("Error converting chatID to int:", err)
			return err
		}

		msg := fmt.Sprintf("Ваш новый пароль: %s", newPsw)
		if err := sendMsgToTg(chatIDInt, msg); err != nil {
			return err
		}
	}

	return users.SetUser(user)
}

func SetUserTgLink(utg *domain.UserTgLink) error {
	return users.SetUserTgLink(utg)
}

func sendMsgToTg(chatIDInt int, msg string) error {
	if err := server.TgClient.SendMessage(chatIDInt, msg); err != nil {
		log.Println("Error sending message to Telegram:", err)
		return err
	}
	log.Println("Message sent to chat ID:", chatIDInt)
	return nil
}
