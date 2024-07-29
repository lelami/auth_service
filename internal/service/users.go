package service

import (
	"authservice/internal/domain"
	"authservice/internal/repository/otpdb"
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
var otps otpdb.DB

func Init(userDB userdb.DB, tokenDB tokendb.DB, otpDB otpdb.DB) {
	users = userDB
	tokens = tokenDB
	otps = otpDB
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

func BlockUserByAdmin(ub *domain.UserBlocked) error {

	user, err := users.GetUser(ub.UserID)
	if err != nil {
		return err
	}

	user.Blocked = ub.Blocked

	return users.SetUser(user)
}

func ChangeRole(ur *domain.UserRole) error {

	user, err := users.GetUser(ur.UserID)
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

	// print to log
	log.Printf("Updated psw %s -> ID %s", newPsw, id.Hex())

	// or sent to TG chat
	chatID, ok := users.CheckExistChatID(id)
	if ok && chatID != nil {
		chatIDInt, err := strconv.Atoi(*chatID)
		if err != nil {
			log.Println("Error converting chatID to int:", err)
			return err
		}

		msg := fmt.Sprintf("Your new password: %s", newPsw)
		if err := sendMsgToTg(chatIDInt, msg); err != nil {
			return err
		}
	}

	return users.SetUser(user)
}

func SetUserTgLink(utg *domain.UserTgLink) error {
	return users.SetUserTgLink(utg)
}

func TgSignIn(l *domain.Login) error {

	userId, ok := users.CheckExistLogin(l.Login)
	if !ok {
		return errors.New("user not found")
	}

	user, err := users.GetUser(*userId)
	if err != nil {
		return err
	}

	code := randOTP()
	otp := &domain.UserOTP{
		UserID:    user.ID,
		Code:      code,
		CreatedAt: time.Now(),
		Expiry:    time.Now().Add(60 * time.Second),
		Used:      false,
	}

	chatID, ok := users.CheckExistChatID(user.ID)
	if ok && chatID != nil {
		chatIDInt, err := strconv.Atoi(*chatID)
		if err != nil {
			log.Println("Error converting chatID to int:", err)
			return err
		}

		msg := fmt.Sprintf("Your one time code: %s", code)
		if err := sendMsgToTg(chatIDInt, msg); err != nil {
			return err
		}
	}
	otps.SetUserOTP(otp)

	return nil
}

func TgCheckOTP(id primitive.ObjectID, uc *domain.Code) (*domain.UserToken, error) {

	userOTP, ok := otps.CheckExistOTP(uc.Code)
	if !ok {
		return nil, errors.New("OTP not found")
	}

	user, err := users.GetUser(id)
	if err != nil {
		return nil, err
	}

	var token string

	if id == userOTP.UserID {
		token = createToken(user.Login)
	} else {
		return nil, errors.New("wrong code")
	}

	if err := otps.MarkOTPAsUsed(uc.Code); err != nil {
		return nil, err
	}

	if err := tokens.SetUserToken(token, id); err != nil {
		return nil, err
	}

	return &domain.UserToken{
		UserId: id,
		Token:  token,
	}, nil
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

func sendMsgToTg(chatIDInt int, msg string) error {
	if err := server.TgClient.SendMessage(chatIDInt, msg); err != nil {
		log.Println("Error sending message to Telegram:", err)
		return err
	}
	log.Println("Message sent to chat ID:", chatIDInt)
	return nil
}
