package httphandler

import (
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SetUserInfoReq struct {
	Name string `json:"name"`
}

type ChangePswReq struct {
	Password string `json:"password"`
}

type BlockUserReq struct {
	ID      primitive.ObjectID `json:"id"`
	Blocked bool               `json:"blocked"`
}

type ChangeRoleReq struct {
	ID   primitive.ObjectID `json:"id"`
	Role string             `json:"role"`
}

type SetUserTgLinkReq struct {
	TgLink string `json:"tg_link"`
}

func (r SetUserInfoReq) IsValid() bool {
	return r.Name != ""
}

// We use regular expression to validate the user tg link.
func (r SetUserTgLinkReq) IsValid() bool {
	reLink := regexp.MustCompile(`^(?:https://t.me/|@)?([a-zA-Z0-9_]{5,32})$`)

	matches := reLink.FindStringSubmatch(r.TgLink)
	if len(matches) != 2 {
		return false
	}

	username := matches[1]

	reUsername := regexp.MustCompile(`^[a-zA-Z0-9_]{5,32}$`)
	return reUsername.MatchString(username)
}

// Retrieving the username for sending messages
func normalizeTgLink(input string) string {
	re := regexp.MustCompile(`^(?:https://t.me/|@)?([a-zA-Z0-9_]{5,32})$`)

	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func (r ChangePswReq) IsValid() bool {
	return r.Password != ""
}

func (r ChangeRoleReq) IsValid() bool {
	if r.Role == "admin" || r.Role == "user" {
		return true
	}
	return false
}
