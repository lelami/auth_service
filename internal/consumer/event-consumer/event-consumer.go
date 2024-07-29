package event_consumer

import (
	"authservice/internal/events"
	"context"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
	ctx       context.Context
}

func New(ctx context.Context, fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		ctx:       ctx,
	}
}
func (c *Consumer) Start() error {
	for {
		select {
		case <-c.ctx.Done():
			log.Println("INFO consumer is shutting down")
			return nil
		default:
			gotEvents, err := c.fetcher.Fetch(c.batchSize)
			if err != nil {
				log.Printf("[ERROR] consumer: %s", err.Error())
				time.Sleep(1 * time.Second)
				continue
			}

			if len(gotEvents) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			if err := c.handleEvents(gotEvents); err != nil {
				log.Print(err)
			}
		}
	}
}

/*
Possible problems and solutions:
1. Loss of events: retrays, return to storage, fallback, confirmation for fetcher
2. Processing the entire batch: stop after the first error, error counter
3. Parallel processing
*/
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
