package cache

import (
	"context"
	"errors"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TelegramAuthCodeCache struct {
	pull map[primitive.ObjectID]int
	mtx  sync.RWMutex
}

const telegramTokenDumpFileName = "telegram_auth_codes.json"

func TelegramAuthCodeCacheInit(ctx context.Context, wg *sync.WaitGroup) (*TelegramAuthCodeCache, error) {
	var c TelegramAuthCodeCache
	c.pull = make(map[primitive.ObjectID]int)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		makeDump(telegramTokenDumpFileName, c.pull)
	}()

	if err := loadFromDump(telegramTokenDumpFileName, &c.pull); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *TelegramAuthCodeCache) GetTelegramAuthCodeByUserId(userID primitive.ObjectID) (int, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	if val, ok := c.pull[userID]; ok {
		return val, nil
	}
	return 0, errors.New("wrong code")
}

func (c *TelegramAuthCodeCache) SetUserTelegramAuthCode(code int, userID primitive.ObjectID) error {
	c.mtx.Lock()
	c.pull[userID] = code
	c.mtx.Unlock()

	return nil
}

func (c *TelegramAuthCodeCache) DeleteUserTelegramAuthCode(userID primitive.ObjectID) error {
	c.mtx.Lock()
	delete(c.pull, userID)
	c.mtx.Unlock()

	return nil
}
