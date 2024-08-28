package otpdb

import "authservice/internal/domain"

type DB interface {
	CheckExistOTP(code string) (*domain.UserOTP, bool)
	SetUserOTP(otp *domain.UserOTP) error
	MarkOTPAsUsed(code string) error
	RemoveExpiredOTPs()
}
