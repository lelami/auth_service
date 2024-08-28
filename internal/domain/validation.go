package domain

func (lp LoginPassword) IsValid() bool {
	if lp.Login == "" || lp.Password == "" {
		return false
	}
	return true
}

func (up UserPassword) IsValid() bool {
	return up.Password == ""
}

func (ui UserInfo) IsValid() bool {
	if ui.ID.IsZero() || ui.Name == "" {
		return false
	}
	return true
}

func (l Login) IsValid() bool {
	return l.Login != ""
}
