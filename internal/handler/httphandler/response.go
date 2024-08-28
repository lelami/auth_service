package httphandler

import "encoding/json"

type HTTPResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (r *HTTPResponse) SetError(err error) {
	r.Error = err.Error()
}

func (r *HTTPResponse) SetData(data interface{}) {
	r.Data = data
}

func (r *HTTPResponse) Marshall() []byte {
	if r.Error == "" {
		r.Success = true
	}
	body, err := json.Marshal(r)
	if err != nil {
		unexpectedResp := &HTTPResponse{
			Error: err.Error(),
		}
		body, _ = json.Marshal(unexpectedResp)
		return body
	}
	return body
}

type ErrorResponse struct {
	Message string `json:"message" example:"login already exists"`
}

type RedirectTgLink struct {
	RedirectURL string `json:"redirect_url"`
}
