package mongo

import (
	"authservice/internal/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"testing"
	"time"
)

func TestMClient_CheckExistLogin(t *testing.T) {
	type fields struct {
		client *mongo.Client
		dbname string
	}
	type args struct {
		login string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *primitive.ObjectID
		want1  bool
	}{
		{
			args: args{login: "four"},
		},
	}

	userDB, err := NewMongoClient("mongodb://admin:admin@localhost:27017/", "auth")
	if err != nil {
		log.Fatalf("ERROR failed to initialize user database: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := userDB.CheckExistLogin(tt.args.login)
			fmt.Println(got, got1)
			/*if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckExistLogin() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CheckExistLogin() got1 = %v, want %v", got1, tt.want1)
			}*/
		})
	}
}

func TestMClient_SetUser(t *testing.T) {
	type fields struct {
		client *mongo.Client
		dbname string
	}
	type args struct {
		user *domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "5",
			args: args{user: &domain.User{
				ID:        primitive.NewObjectID(),
				Login:     "5",
				Password:  "",
				Name:      "John",
				Role:      "user",
				Skills:    []string{"go", "linux", "git"},
				Education: nil,
				Created:   time.Now(),
				Updated:   time.Time{},
			}},
			wantErr: false,
		},
		{
			name: "5",
			args: args{user: &domain.User{
				ID:        primitive.NewObjectID(),
				Login:     "5",
				Password:  "",
				Name:      "John",
				Role:      "user",
				Skills:    []string{"go", "linux", "git"},
				Education: nil,
				Created:   time.Now(),
				Updated:   time.Time{},
			}},
			wantErr: true,
		},
		{
			name: "6",
			args: args{user: &domain.User{
				ID:        primitive.NewObjectID(),
				Login:     "7",
				Password:  "",
				Name:      "John",
				Role:      "user",
				Skills:    []string{"go", "git"},
				Education: nil,
				Created:   time.Now(),
				Updated:   time.Time{},
			}},
			wantErr: false,
		},
	}

	userDB, err := NewMongoClient("mongodb://admin:admin@localhost:27017/", "auth")
	if err != nil {
		log.Fatalf("ERROR failed to initialize user database: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := userDB.SetUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("SetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMClient_GetUserByCreated(t *testing.T) {
	type fields struct {
		client *mongo.Client
		dbname string
	}
	type args struct {
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.User
		wantErr bool
	}{
		{
			name: "",
			args: args{
				from: time.Now().Add(-20 * time.Minute),
				to:   time.Now(),
			},
		},
	}

	userDB, err := NewMongoClient("mongodb://admin:admin@localhost:27017/", "auth")
	if err != nil {
		log.Fatalf("ERROR failed to initialize user database: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := userDB.GetUserByCreated(tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByCreated() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}
