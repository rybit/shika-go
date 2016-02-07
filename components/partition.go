package components

import (
	"encoding/json"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Partition this is an actual data container. They can be proxies to a remote
// machine or just a local machine
type partition interface {
	Write(payload string) error
	Subscribe(listener chan Message)
	SubscribersCount() int
	Close()
}

// -------------------------------------------------------------------------------------------------
// REMOTE partition
// -------------------------------------------------------------------------------------------------

type remotePartition struct {
	host string
	port int
}

func (rp remotePartition) Write(payload string) error {
	return nil // TODO
}

func (rp remotePartition) Subscribe(listener chan Message) {
	// TODO
}

func (rp remotePartition) SubscribersCount() int {
	// TODO
	return 0
}

func (rp remotePartition) Close() {
	// TODO
}

// -------------------------------------------------------------------------------------------------
// LOCAL partition
// -------------------------------------------------------------------------------------------------

type localPartition struct {
	topicName    string
	dataFile     *os.File
	currentMsgID uint64
	listeners    []chan Message

	sync.Mutex
}

func newLocalPartition(filename string, topic string) (partition, error) {
	// outputFile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	outputFile, err := os.Create(filename) // b/c we don't have a way to read from the file
	if err != nil {
		return nil, err
	}

	lp := new(localPartition)
	lp.dataFile = outputFile
	lp.topicName = topic

	return lp, err
}

func (l *localPartition) SubscribersCount() int {
	return len(l.listeners)
}

func (l *localPartition) Close() {
	l.dataFile.Close()
}

func (l *localPartition) Write(payload string) error {
	thisMsgID := atomic.AddUint64(&l.currentMsgID, 1)
	m := Message{
		ID:        thisMsgID,
		Payload:   payload,
		Timestamp: time.Now(),
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
	l.Lock()
	defer l.Unlock()

	l.listeners = append(l.listeners, listener)
}

func informListener(listener chan Message, m Message) {
	listener <- m
}
