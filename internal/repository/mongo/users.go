package mongo

import (
	"authservice/internal/domain"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const CollUsers = "users"

func (c *MClient) CheckExistLogin(login string) (*primitive.ObjectID, bool) {
	coll := c.client.Database(c.dbname).Collection(CollUsers)

	filter := bson.M{"login": login}

	var user domain.User

	if err := coll.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return nil, false
	}
	return &user.ID, true
}

func (c *MClient) GetUser(id primitive.ObjectID) (*domain.User, error) {
	coll := c.client.Database(c.dbname).Collection(CollUsers)

	filter := bson.M{"_id": id}

	var user domain.User

	if err := coll.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *MClient) SetUser(user *domain.User) error {

	user.Updated = time.Now()
	filter := bson.D{{"_id", user.ID}}
	set := bson.D{{"$set", user}}

	coll := c.client.Database(c.dbname).Collection(CollUsers)
	opts := options.FindOneAndUpdate().SetUpsert(true)
	res := coll.FindOneAndUpdate(context.TODO(), filter, set, opts)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return nil
	}
	fmt.Println(res.Err())
	return res.Err()
}

func (c *MClient) GetUserByCreated(from, to time.Time) ([]domain.User, error) {
	coll := c.client.Database(c.dbname).Collection(CollUsers)

	filter := bson.M{
		"$and": bson.A{
			bson.M{"created": bson.M{"$gte": from}},
			bson.M{"created": bson.M{"$lte": to}},
		},
	}
	// доп примеры поиска
	/*	filter := bson.M{"skills": bson.M{"$all": bson.A{"go", "linux"}}}*/
	//{ "skills": { $all: ["Go", "Git"] }

	var users []domain.User
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err := cur.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return users, nil
}
