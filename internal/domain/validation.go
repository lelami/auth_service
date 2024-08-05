package domain

func (lp LoginPassword) IsValid() bool {
	if lp.Login == "" || lp.Password == "" {
		return false
	}
	return true
}
func (lp LoginCode) IsValid() bool {
	if lp.Login == "" || lp.Code == "" {
		return false
	}
	return true
}
func (up UserPassword) IsValid() bool {
	if up.Password == "" {
		return false
	}
	return true
}

func (ui UserInfo) IsValid() bool {
	if ui.ID.IsZero() || ui.Name == "" {
		return false
	}
	return true
}
