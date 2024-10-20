package kafka

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestClient_Consume(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	cl, err := Init(ctx, "localhost:9094", "test2", "test3")
	if err != nil {
		t.Errorf("Init() error = %v", err)
	}

	printMsg := func(topic string, groupId string, msg []byte) {
		fmt.Printf("Message on %s groupid %s: %s\n", topic, groupId, string(msg))
	}

	go func() {
		if err := cl.Consume(ctx, "6", printMsg, "test2"); err != nil {
			t.Errorf("Consume() error = %v", err)
		}
	}()

	go func() {
		if err := cl.Consume(ctx, "8", printMsg, "test2"); err != nil {
			t.Errorf("Consume() error = %v", err)
		}
	}()

	time.Sleep(10 * time.Second)
	cancel()
	time.Sleep(10 * time.Second)
}

func TestClient_SendMessage(t *testing.T) {

	type args struct {
		topic   string
		message string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "send to test2",
			args: args{
				topic:   "test2",
				message: "test2 message at lecture 2",
			},
			wantErr: false,
		},
		{
			name: "send to test3",
			args: args{
				topic:   "test3",
				message: "test3 message at lecture",
			},
			wantErr: false,
		},
	}

	kfkClient, err := Init(context.Background(), "localhost:9094", "test2", "test3")
	if err != nil {
		t.Errorf("Init() error = %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := kfkClient.SendMessage(tt.args.topic, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	kfkClient.ProducerWg.Wait()
}

func TestInit(t *testing.T) {
	type args struct {
		url    string
		topics []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				url:    "localhost:9094",
				topics: []string{"test2", "test3"},
			},
			want: &Client{
				servers:       "localhost:9094",
				produceTopics: map[string]struct{}{"test2": {}, "test3": {}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Init(context.Background(), tt.args.url, tt.args.topics...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.servers, tt.want.servers) {
				t.Errorf("Init() got = %v, want %v", got.servers, tt.want.servers)
			}
			if !reflect.DeepEqual(got.produceTopics, tt.want.produceTopics) {
				t.Errorf("Init() got = %v, want %v", got.produceTopics, tt.want.produceTopics)
			}
		})
	}
}
