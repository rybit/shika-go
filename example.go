package main

import (
	"fmt"
	"sync"

	c "github.com/squirrely/shikago/components"
)

func main() {
	incoming := make(chan c.Message)

	part, err := c.NewLocalPartition("/tmp/shikago-test.jsonl")
	if err != nil {
		panic(err)
	}
	defer part.Close()
	part.Subscribe(incoming)

	go func() {
		fmt.Println("Starting to consume")
		for msg := range incoming {
			fmt.Printf("%d: %s\n", msg.Id, msg.Payload)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go write(part, &wg, "s1", 40)
	go write(part, &wg, "s2", 50)

	wg.Wait()
	fmt.Println("Finished writing")
}

func write(part c.Partition, wg *sync.WaitGroup, senderID string, max int) {
	fmt.Println("Starting to produce: " + senderID)
	for i := 0; i < max; i++ {
		part.Write(fmt.Sprintf("%s -- %d", senderId, i))
	}
	wg.Done()
}
