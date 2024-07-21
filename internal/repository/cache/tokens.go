package cache

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type TokenCache struct {
	pull map[string]primitive.ObjectID
	mtx  sync.RWMutex
}

const tokenDumpFileName = "tokens.json"

func TokenCacheInit(ctx context.Context, wg *sync.WaitGroup) (*TokenCache, error) {
	var c TokenCache
	c.pull = make(map[string]primitive.ObjectID)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		makeDump(tokenDumpFileName, c.pull)
	}()

	if err := loadFromDump(tokenDumpFileName, &c.pull); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *TokenCache) GetUserByToken(token string) (*primitive.ObjectID, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	if val, ok := c.pull[token]; ok {
		return &val, nil
	}
	return nil, errors.New("token not found")
}

func (c *TokenCache) SetUserToken(token string, userID primitive.ObjectID) error {

	c.mtx.Lock()
	c.pull[token] = userID
	c.mtx.Unlock()

	return nil
}
