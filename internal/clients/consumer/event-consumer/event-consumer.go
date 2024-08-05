package event_consumer

import (
	"authservice/internal/clients/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	partSize  int
}

func New(fetcher events.Fetcher, processor events.Processor, partSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		partSize:  partSize,
	}
}
func (consumer *Consumer) Start() {
	for {
		gotEvents, err := consumer.fetcher.Fetch(consumer.partSize)
		if err != nil {
			log.Printf("Error fetching events: %s", err)
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := consumer.HandleEvents(gotEvents); err != nil {
			log.Printf("Error handling events: %s", err)
			continue
		}

	}
}
func (consumer *Consumer) HandleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("Processing event: %s", event.Text)
		if err := consumer.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err)
			continue
		}
	}
	return nil
}
