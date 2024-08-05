package httphandler

type SetUserInfoReq struct {
	Name string `json:"name"`
}
type TelegramName struct {
	Name string `json:"name"`
}
type ChangePswReq struct {
	Password string `json:"password"`
}

func (r SetUserInfoReq) IsValid() bool {
	return r.Name != ""
}

func (r ChangePswReq) IsValid() bool {
	return r.Password != ""
}

func (t TelegramName) IsValid() bool {
	return t.Name != ""
}
