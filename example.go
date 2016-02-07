package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	c "github.com/squirrely/shikago/components"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %s <topic name> <config file>\n", os.Args[0])
		os.Exit(1)
	}
	fileBytes, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println("Failed to read in file", os.Args[2])
		panic(err)
	}

	var config c.Configuration
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		fmt.Println("Failed to unmarshal file", os.Args[2])
		panic(err)
	}

	topic := os.Args[1]

	node := c.NewNode(&config)
	defer node.Shutdown()

	all := node.SubscribeToAll(topic)

	go consume("all consumer", all)

	// b/c we know how it hands out the conusmers to the one with the lowest number of
	// consumers, we can reliably know that this will cause a consumer per partition
	for i := 0; int32(i) < config.DefaultPartitionCount; i++ {
		incoming := node.Subscribe(topic)
		consumerName := "consumer-" + strconv.Itoa(i)

		go consume(consumerName, incoming)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		time.Sleep(10 * time.Millisecond)
		fmt.Print("Write to partition: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		text = strings.TrimSpace(text)
		node.Write(topic, text)
	}
}

func consume(id string, incoming <-chan c.Message) {
	fmt.Println("Starting", id)
	for msg := range incoming {
		fmt.Printf("%s - msg %d: %s\n", id, msg.ID, msg.Payload)
	}
}

//
// func write(wg *sync.WaitGroup, node *c.Node, senderID string, max int) {
// 	fmt.Println("Starting to produce: " + senderID)
//
// 	for i := 0; i < max; i++ {
// 		payload := fmt.Sprintf("sender %s:%d", senderID, i)
// 		node.Write(TOPIC, payload)
// 	}
//
// 	wg.Done()
// }
