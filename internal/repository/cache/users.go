package cache

import (
	"authservice/internal/domain"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type UserCache struct {
	userPool  map[primitive.ObjectID]*domain.User
	loginPool map[string]primitive.ObjectID
	mtx       sync.RWMutex
}

const userDumpFileName = "users.json"

func UserCacheInit(ctx context.Context, wg *sync.WaitGroup) (*UserCache, error) {

	var c UserCache
	c.userPool = make(map[primitive.ObjectID]*domain.User)
	c.loginPool = make(map[string]primitive.ObjectID)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		makeDump(userDumpFileName, c.userPool)
	}()

	if err := loadFromDump(userDumpFileName, &c.userPool); err != nil {
		return nil, err
	}

	for _, user := range c.userPool {
		c.loginPool[user.Login] = user.ID
	}

	return &c, nil
}

/*
func (c *UserCache) CheckExistLogin(login string) (*primitive.ObjectID, bool) {

		c.mtx.RLock()
		id, ok := c.loginPool[login]
		c.mtx.RUnlock()

		return &id, ok
	}
*/
func (c *UserCache) CheckExistLogin(login string) (*primitive.ObjectID, bool) {

	c.mtx.RLock()
	id, ok := c.loginPool[login]
	c.mtx.RUnlock()

	return &id, ok
}

func (c *UserCache) GetUser(id primitive.ObjectID) (*domain.User, error) {

	c.mtx.RLock()
	user, ok := c.userPool[id]
	c.mtx.RUnlock()

	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (c *UserCache) SetUser(newUserInfo *domain.User) error {

	c.mtx.Lock()
	c.userPool[newUserInfo.ID] = newUserInfo
	c.loginPool[newUserInfo.Login] = newUserInfo.ID
	c.mtx.Unlock()

	return nil
}
