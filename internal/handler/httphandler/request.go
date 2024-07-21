package httphandler

type SetUserInfoReq struct {
	Name string `json:"name"`
}

func (r SetUserInfoReq) IsValid() bool {
	return r.Name != ""
}

type ChangePswReq struct {
	Password string `json:"password"`
}

func (r ChangePswReq) IsValid() bool {
	return r.Password != ""
}
