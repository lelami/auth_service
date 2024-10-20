package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"sync"
	"time"
)

type Client struct {
	servers       string
	cfg           *kafka.ConfigMap
	serviceCtx    context.Context
	producer      *kafka.Producer
	ProducerWg    *sync.WaitGroup
	produceTopics map[string]struct{}
}

func Init(ctx context.Context, url string, topics ...string) (*Client, error) {

	if len(topics) == 0 {
		return &Client{servers: url}, nil
	}

	cfg := &kafka.ConfigMap{
		"bootstrap.servers": url,
	}

	wg := &sync.WaitGroup{}
	client := &Client{
		servers:    url,
		cfg:        cfg,
		serviceCtx: ctx,
		ProducerWg: wg,
	}

	if err := client.createTopics(topics); err != nil {
		return nil, err
	}

	if err := client.createProducer(); err != nil {
		return nil, err
	}

	return client, nil
}

func (cl *Client) createTopics(topics []string) error {
	p, err := kafka.NewAdminClient(cl.cfg)
	if err != nil {
		return err
	}
	defer p.Close()

	topicsScep := make([]kafka.TopicSpecification, 0, len(topics))
	for _, topic := range topics {
		topicsScep = append(topicsScep, kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
			Config: map[string]string{
				"retention.ms":      "86400000", // 1 день
				"max.message.bytes": "1048576",  // 1 МБ
			},
		})
	}

	results, err := p.CreateTopics(
		context.Background(),
		topicsScep,
		kafka.SetAdminOperationTimeout(50000*time.Millisecond))

	if err != nil {
		return err
	}
	topicMap := make(map[string]struct{})

	topicsErr := ""
	for _, result := range results {
		switch result.Error.Code() {
		case kafka.ErrTopicAlreadyExists:
			topicMap[result.Topic] = struct{}{}
			fmt.Printf("Topic %s is already exist\n", result.Topic)
		case kafka.ErrNoError:
			topicMap[result.Topic] = struct{}{}
			fmt.Printf("Topic %s created\n", result.Topic)
		default:
			topicsErr += result.Error.Error() + ", "
		}
	}

	if topicsErr != "" {
		return errors.New(topicsErr)
	}

	cl.produceTopics = topicMap
	return nil
}

func (cl *Client) createProducer() error {
	p, err := kafka.NewProducer(cl.cfg)
	if err != nil {
		return err
	}
	cl.producer = p

	// Запуск горутины для обработки событий доставки
	go func() {
		for {
			select {
			case e := <-p.Events():
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
					} else {
						fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
					}
					cl.ProducerWg.Done()
					log.Println("wg.Done for", ev.TopicPartition)
				}
			case <-cl.serviceCtx.Done():
				cl.ProducerWg.Wait()
				cl.CloseProducer()
				return
			}
		}
	}()

	return nil
}

func (cl *Client) CloseProducer() {
	cl.producer.Flush(15 * 1000)
	cl.producer.Close()
}

func (cl *Client) SendMessage(topic string, message string) error {
	if _, ok := cl.produceTopics[topic]; !ok {
		return fmt.Errorf("topic %s is not registered", topic)
	}
	cl.ProducerWg.Add(1)
	err := cl.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	return err
}

func (cl *Client) Consume(ctx context.Context, groupId string, handler func(string, string, []byte), topics ...string) error {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cl.servers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return err
	}
	rebalanceCb := func(c *kafka.Consumer, event kafka.Event) error {
		switch e := event.(type) {
		case kafka.AssignedPartitions:
			log.Printf("Partitions assigned: %v\n", e.Partitions)
			c.Assign(e.Partitions)
		case kafka.RevokedPartitions:
			log.Printf("Partitions revoked: %v\n", e.Partitions)
			c.Unassign()
		}
		return nil
	}
	if err := c.SubscribeTopics(topics, rebalanceCb); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled, stopping consumer")
			err = c.Close()
			return err
		default:
			msg, err := c.ReadMessage(-1)
			if err == nil {
				handler(*msg.TopicPartition.Topic, groupId, msg.Value)
			} else {
				fmt.Printf("Consumer error groupid %s: %v (%v)\n", groupId, err, msg)
			}
		}
	}
}
