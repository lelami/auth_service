package httphandler

import (
	"authservice/internal/domain"
	"authservice/internal/service"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminGetUserInfo godoc
// @Tags AdminAuth
// @Summary Получение личных данных пользователя админом
// @Param Authorization header string true "токен админа"
// @Param user_id query string true "id пользователя"
// @Success 200 {object} domain.User
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /admin/get_user_info [get]
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

// AdminBlockUser godoc
// @Tags AdminAuth
// @Summary Блокировка пользователя админом
// @Param Authorization header string true "токен админа"
// @Param request body BlockUserReq true "id пользователя и значение блокировки"
// @Success 200
// @Failure 404
// @Failure 409
// @Failure 422
// @Failure 500
// @Router /admin/block_user [post]
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
		UserID:  input.ID,
		Blocked: input.Blocked,
	})
	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		respBody.SetError(err)
		return
	}
}

// AdminChangeRole godoc
// @Tags AdminAuth
// @Summary Смена роли пользователя админом
// @Param Authorization header string true "токен админа"
// @Param request body ChangeRoleReq true "id пользователя и роль"
// @Success 200
// @Failure 404
// @Failure 409
// @Failure 422
// @Failure 500
// @Router /admin/change_role [post]
func AdminChangeRole(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input ChangeRoleReq

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

	err := service.ChangeRole(&domain.UserRole{
		UserID: input.ID,
		Role:   input.Role,
	})
	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		respBody.SetError(err)
		return
	}
}
