package cache

import (
	"authservice/internal/domain"
	"context"
	"errors"
	"sync"
	"time"
)

type OTPCache struct {
	otpPull map[string]*domain.UserOTP
	mtx     sync.RWMutex
}

const otpDumpFileName = "otps.json"

// OTPCacheInit initializes the cache for one-time codes and loads data from the dump.
func OTPCacheInit(ctx context.Context, wg *sync.WaitGroup) (*OTPCache, error) {
	var c OTPCache
	c.otpPull = make(map[string]*domain.UserOTP)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		makeDump(otpDumpFileName, c.otpPull)
	}()

	if err := loadFromDump(otpDumpFileName, &c.otpPull); err != nil {
		return nil, err
	}

	return &c, nil
}

// CheckExistOTP проверяет существование одноразового кода и его статус.
func (c *OTPCache) CheckExistOTP(code string) (*domain.UserOTP, bool) {
	c.mtx.RLock()
	otp, ok := c.otpPull[code]
	c.mtx.RUnlock()

	return otp, ok
}

// SetUserOTP adds or updates a one-time code in the cache.
func (c *OTPCache) SetUserOTP(otp *domain.UserOTP) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.otpPull[otp.Code] = otp

	return nil
}

// MarkOTPAsUsed marks the one-time code as used.
func (c *OTPCache) MarkOTPAsUsed(code string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	otp, ok := c.otpPull[code]
	if !ok {
		return errors.New("OTP not found")
	}

	if otp.Used {
		return errors.New("OTP already used")
	}

	otp.Used = true

	return nil
}

// RemoveExpiredOTPs removes expired one-time codes.
// Ideally, run a daemon to remove the overdue.
func (c *OTPCache) RemoveExpiredOTPs() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for code, otp := range c.otpPull {
		if time.Now().After(otp.Expiry) {
			delete(c.otpPull, code)
		}
	}
}
