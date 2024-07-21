package httphandler

import (
	"authservice/internal/domain"
	"authservice/internal/service"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/sign_up", CORS(LogUser(http.HandlerFunc(SignUp))))
	router.Handle("/sign_in", CORS(LogUser(http.HandlerFunc(SignIn))))
	router.Handle("/get_user_info", CORS(Auth(LogUser(http.HandlerFunc(GetUserInfo)))))
	router.Handle("/set_user_info", CORS(Auth(LogUser(http.HandlerFunc(SetUserInfo)))))
	router.Handle("/change_psw", CORS(Auth(LogUser(http.HandlerFunc(ChangePsw)))))

	return router
}

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

	info, err := service.GetUserInfo(userID)
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
