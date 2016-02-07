package main

import (
	"fmt"
	"sync"

	c "github.com/squirrely/shikago/components"
)

func main() {
	incoming1 := make(chan c.Message)
	incoming2 := make(chan c.Message)
	incoming3 := make(chan c.Message)

	topic, err := c.NewTopic("/tmp/shikago", "test-data", 2)
	if err != nil {
		panic(err)
	}
	defer topic.Close()

	topic.Subscribe(incoming1)
	topic.Subscribe(incoming2)
	topic.SubscribeToAll(incoming3)

	go consume("01", incoming1)
	go consume("02", incoming2)
	go consume("03", incoming3)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go write(topic, &wg, "s1", 10)
	go write(topic, &wg, "s2", 5)

	wg.Wait()
	fmt.Println("Finished writing")
}

func consume(id string, incoming chan c.Message) {
	fmt.Println("Starting to consume - ", id)
	for msg := range incoming {
		fmt.Printf("consumer %s - %d: %s\n", id, msg.ID, msg.Payload)
	}
}

func write(topic *c.Topic, wg *sync.WaitGroup, senderID string, max int) {
	fmt.Println("Starting to produce: " + senderID)
	for i := 0; i < max; i++ {
		topic.Write(fmt.Sprintf("sender %s -- %d", senderID, i))
	}
	wg.Done()
}
