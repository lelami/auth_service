package cache

import (
	"authservice/internal/domain"
	"context"
	"errors"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserCache struct {
	userPull  map[primitive.ObjectID]*domain.User
	loginPull map[string]primitive.ObjectID
	tgPull    map[string]primitive.ObjectID
	mtx       sync.RWMutex
}

const userDumpFileName = "users.json"

func UserCacheInit(ctx context.Context, wg *sync.WaitGroup) (*UserCache, error) {

	var c UserCache
	c.userPull = make(map[primitive.ObjectID]*domain.User)
	c.loginPull = make(map[string]primitive.ObjectID)
	c.tgPull = make(map[string]primitive.ObjectID)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		makeDump(userDumpFileName, c.userPull)
	}()

	if err := loadFromDump(userDumpFileName, &c.userPull); err != nil {
		return nil, err
	}

	for _, user := range c.userPull {
		c.loginPull[user.Login] = user.ID
		c.tgPull[user.TgLink] = user.ID
	}

	return &c, nil
}

func (c *UserCache) CheckExistLogin(login string) (*primitive.ObjectID, bool) {

	c.mtx.RLock()
	id, ok := c.loginPull[login]
	c.mtx.RUnlock()

	return &id, ok
}

func (c *UserCache) GetUser(id primitive.ObjectID) (*domain.User, error) {

	c.mtx.RLock()
	user, ok := c.userPull[id]
	c.mtx.RUnlock()

	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (c *UserCache) SetUser(newUserInfo *domain.User) error {

	c.mtx.Lock()
	c.userPull[newUserInfo.ID] = newUserInfo
	c.loginPull[newUserInfo.Login] = newUserInfo.ID
	c.mtx.Unlock()

	return nil
}

func (c *UserCache) SetUserTgLink(utg *domain.UserTgLink) error {
	user, err := c.GetUser(utg.UserID)
	if err != nil {
		return err
	}

	c.mtx.Lock()
	c.tgPull[utg.TgLink] = utg.UserID
	c.mtx.Unlock()
	log.Printf("Updated tgPull with link %s -> ID %s", utg.TgLink, utg.UserID.Hex())

	user.TgLink = utg.TgLink

	return c.SetUser(user)
}

func (c *UserCache) GetUserByTgLink(tgLink string) (*primitive.ObjectID, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	if val, ok := c.tgPull[tgLink]; ok {
		return &val, nil
	}
	return nil, errors.New("tg link not found")
}

func (c *UserCache) SetUserChatID(chatID *domain.UserChatID) error {
	userID, err := c.GetUserByTgLink(chatID.TgLink)
	if err != nil {
		return err
	}

	user, err := c.GetUser(*userID)
	if err != nil {
		return err
	}

	c.mtx.Lock()
	user.ChatID = chatID.ChatID
	c.mtx.Unlock()

	return c.SetUser(user)
}

func (c *UserCache) CheckExistChatID(id primitive.ObjectID) (*string, bool) {

	c.mtx.RLock()
	user, ok := c.userPull[id]
	c.mtx.RUnlock()

	if !ok {
		return nil, false
	}

	chatIDStr := user.ChatID
	return &chatIDStr, true
}
