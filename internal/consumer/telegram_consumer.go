package consumer

import (
	event_consumer "authservice/internal/clients/consumer/event-consumer"
	"authservice/internal/clients/events/telegram"
	tgClient "authservice/internal/clients/telegram"
)

const (
	partSize = 100
)

var eventProc *telegram.FetchProcessor
var consumer *event_consumer.Consumer

func Init(client *tgClient.Client) {
	eventProc = telegram.New(client)
	consumer = event_consumer.New(eventProc, eventProc, partSize)

}
func StartTelegramConsumer() {

	go func() {
		consumer.Start()
	}()

}
