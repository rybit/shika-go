package main

import (
	"fmt"
	"sync"
	"time"

	c "github.com/squirrely/shikago/components"
)

// TOPIC ==> junk
const TOPIC = "junk-topic"

func main() {
	config := c.Configuration{
		DefaultPartitionSize: 2,
		DataDirectory:        "/tmp/shikago",
	}

	node := c.NewNode(&config)
	incoming1 := node.Subscribe(TOPIC)
	incoming2 := node.Subscribe(TOPIC)
	all := node.SubscribeToAll(TOPIC)

	go consume("first consumer", incoming1)
	go consume("second consumer", incoming2)
	go consume("all consumer", all)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go write(&wg, node, "s1", 10)
	go write(&wg, node, "s2", 10)

	wg.Wait()
	fmt.Println("Finished writing")

	time.Sleep(time.Second * 2)
}

func consume(id string, incoming <-chan c.Message) {
	fmt.Println("Starting to consume -", id)
	for msg := range incoming {
		fmt.Printf("%s - msg %d: %s\n", id, msg.ID, msg.Payload)
	}
}

func write(wg *sync.WaitGroup, node *c.Node, senderID string, max int) {
	fmt.Println("Starting to produce: " + senderID)

	for i := 0; i < max; i++ {
		payload := fmt.Sprintf("sender %s:%d", senderID, i)
		node.Write(TOPIC, payload)
	}

	wg.Done()
}
