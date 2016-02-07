package components

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	s "strings"
	"sync"
	"sync/atomic"
)

// Message the container for persisted data. The ID field is considered
// unique PER partition.
type Message struct {
	ID      uint64 `json:"id"`
	Payload string `json:"payload"`
}

// Topic a container for all of the partitions and the strategy
// which defines to which partition to route the data
type Topic struct {
	Name           string
	PartitionCount int32

	strategy RoutingStrategy
	parts    []partition
}

// Partition this is an actual data container. They can be proxies to a remote
// machine or just a local machine
type partition interface {
	Write(payload string) error
	Subscribe(listener chan Message)
	GetNumberOfSubscribers() int
	Close()
}

// -------------------------------------------------------------------------------------------------
// Topic
// -------------------------------------------------------------------------------------------------

// NewTopic create a topic with the round robin strategy. The underlying partition files
// will be contained in the directory specified as "topic_name_1.jsonl", "topic_name_2.jsonl", etc
//
// TODO: think more about how we will be doing global configuration
//
func NewTopic(dataDirPath string, name string, partCount int32) (*Topic, error) {
	// make the local partitions
	parts := make([]partition, partCount)

	for i := range parts {
		partName := fmt.Sprintf("%s_%d.jsonl", name, partCount)
		partName = s.ToLower(s.TrimSpace(s.Replace(partName, " ", "_", -1)))

		part, err := newLocalPartition(filepath.Join(dataDirPath, partName))
		if err != nil {
			log.Println("Failed to create partition ", partName)
			return nil, err
		}

		parts[i] = part
	}

	topic := Topic{
		Name:           name,
		PartitionCount: partCount,
		parts:          parts,
		strategy:       NewRoundRobinStrategy(partCount),
	}
	return &topic, nil
}

// Write will decide which partition to route the request to and then issue the write
func (t Topic) Write(payload string) error {
	partitionIndex := t.strategy.WhichPartition(payload)

	return t.parts[partitionIndex].Write(payload)
}

// Subscribe finds the partition with the least subscribers and then delegates to that partition
func (t Topic) Subscribe(listener chan Message) {
	var part partition
	for _, p := range t.parts {
		if part == nil {
			part = p
		} else {
			if p.GetNumberOfSubscribers() > part.GetNumberOfSubscribers() {
				part = p
			}
		}
	}

	part.Subscribe(listener)
}

// SubscribeToAll will bind that listener to each of the underlying partitions
func (t Topic) SubscribeToAll(listener chan Message) {
	for _, p := range t.parts {
		p.Subscribe(listener)
	}
}

// Close closes each partition gracefully
func (t Topic) Close() {
	for _, part := range t.parts {
		part.Close()
	}
}

// -------------------------------------------------------------------------------------------------
// LOCAL partition
// -------------------------------------------------------------------------------------------------

type localPartition struct {
	dataFile      *os.File
	currentMsgID  uint64
	listenerMutex *sync.Mutex
	listeners     []chan Message
}

func newLocalPartition(filepath string) (partition, error) {
	outputFile, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	lp := new(localPartition)
	lp.dataFile = outputFile

	return lp, err
}

func (l localPartition) GetNumberOfSubscribers() int {
	return len(l.listeners)
}

func (l localPartition) Close() {
	l.dataFile.Close()
}

func (l *localPartition) Write(payload string) error {
	thisMsgID := atomic.AddUint64(&l.currentMsgID, 1)
	m := Message{
		ID:      thisMsgID,
		Payload: payload,
	}

	// write it to file first
	jsonPayload, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// for now write it as a jsonl file - easier debugging
	_, err = l.dataFile.WriteString(string(jsonPayload) + "\n")
	if err != nil {
		return err
	}

	// tell everyone we about the data
	for _, listener := range l.listeners {
		go informListener(listener, m)
	}

	return nil
}

func (l *localPartition) Subscribe(listener chan Message) {
	// l.listenerMutex.Lock()
	// defer l.listenerMutex.Unlock()

	l.listeners = append(l.listeners, listener)
}

func informListener(listener chan Message, m Message) {
	listener <- m
}
