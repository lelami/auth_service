package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type RMQClient struct {
	Conn      *amqp.Connection
	publishCH *amqp.Channel
}

// Init Устанавливает соединение с RabbitMQ
func Init(url string, queue string) (*RMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	c := &RMQClient{
		Conn: conn,
	}
	ch, err := c.Conn.Channel()
	if err != nil {
		return nil, err
	}
	c.publishCH = ch

	_, err = ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *RMQClient) Close() error {
	return c.Conn.Close()
}

// DeclareQueue Объявляет очередь
func (c *RMQClient) DeclareQueue(name string) (amqp.Queue, error) {
	ch, err := c.Conn.Channel()
	if err != nil {
		return amqp.Queue{}, err
	}
	defer ch.Close()

	/*	arguments := amqp.Table{
		"x-message-ttl":             int32(60000),  // TTL сообщений в миллисекундах
		"x-expires":                 int32(300000), // TTL очереди в миллисекундах
		"x-max-length":              int32(1000),   // Максимальное количество сообщений в очереди
		"x-dead-letter-exchange":    "dead_letter_exchange",
		"x-dead-letter-routing-key": "dead_letter_queue",
	}*/
	return ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// DeclareExchange Объявляет обменник
func (c *RMQClient) DeclareExchange(name, kind string) error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	arguments := amqp.Table{
		"alternate-exchange": "alternate_exchange", // Альтернативный обменник
	}
	return ch.ExchangeDeclare(
		name,      // name
		kind,      // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		arguments, // arguments
	)
}

// BindQueueToExchange Привязывает очередь к обменнику
func (c *RMQClient) BindQueueToExchange(queueName, exchangeName, routingKey string) error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
}

// Publish Отправляет сообщение в обменник
func (c *RMQClient) Publish(exchange, routingKey string, body []byte) error {
	return c.publishCH.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

// PublishToQueue Отправляет сообщение в очередь
func (c *RMQClient) PublishToQueue(queueName string, body []byte) error {
	// Отправляем сообщение в очередь
	return c.publishCH.Publish(
		"",        // exchange (пустая строка означает отправку напрямую в очередь)
		queueName, // routing key (имя очереди)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

// Consume Получает сообщения из очереди
func (c *RMQClient) Consume(queueName string, handle func(msg interface{}) error) error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}

	// Регистрируем потребителя с ручным подтверждением
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (ручное подтверждение)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	// Обрабатываем сообщения
	go func() {
		for msg := range msgs {

			err := handle(msg)
			if err != nil {
				msg.Nack(false, false)
			}
			// Ручное подтверждение сообщения
			err = msg.Ack(false)
			if err != nil {
				log.Printf("Failed to acknowledge message: %v", err)
			} else {
				log.Println("Message acknowledged")
			}
		}
	}()

	return nil
}

func (c *RMQClient) DeclareStream(streamName string) error {

	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	// Дополнительные аргументы для стрима
	args := amqp.Table{
		"x-max-length-bytes":              20000000000, // 20 GB
		"x-max-age":                       "7D",        // 7 дней
		"x-stream-max-segment-size-bytes": 100000000,   // 100 MB
	}
	// Создаем стрим
	return ch.ExchangeDeclare(
		streamName, // name
		"stream",   // kind
		true,       // durable
		false,      // autoDelete
		false,      // internal
		false,      // noWait
		args,       // arguments
	)
}

func (c *RMQClient) PublishToStream(streamName string) error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	// Публикация сообщения с заголовком x-stream-filter-value
	return ch.Publish(
		streamName,
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Message from California"),
			Headers: amqp.Table{
				"x-stream-filter-value": "california", // установка значения фильтра
			},
		},
	)
}

// Consume Получает сообщения из очереди
func (c *RMQClient) ConsumeByFilter(queueName, filter string) error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}

	// Регистрируем потребителя с ручным подтверждением
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (ручное подтверждение)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		amqp.Table{
			"x-stream-filter": "california", // установка фильтра
		}, // args
	)
	if err != nil {
		return err
	}

	// Обрабатываем сообщения
	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			headers := msg.Headers
			filterValue, ok := headers["x-stream-filter-value"].(string)
			if ok && filterValue == "california" {
				log.Printf("Received a message: %s", msg.Body)
				// Обработка сообщения
				// ...
			}
			// Ручное подтверждение сообщения
			err := msg.Ack(false)
			if err != nil {
				log.Printf("Failed to acknowledge message: %v", err)
			} else {
				log.Println("Message acknowledged")
			}
		}
	}()

	return nil
}
