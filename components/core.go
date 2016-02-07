package components

import (
	"encoding/json"
	"os"
	"sync"
	"sync/atomic"
)

type Message struct {
	Id      uint64 `json:"id"`
	Payload string `json:"payload"`
}

type Topic struct {
	parts []Partition
}

type Partition interface {
	Write(payload string) error
	Subscribe(listener chan Message)
	Close()
}

// ------------------------------------------------------------------------------------------------
// LOCAL partition
// ------------------------------------------------------------------------------------------------

type localPartition struct {
	dataFile      *os.File
	currentMsgId  uint64
	listenerMutex sync.Mutex
	listeners     []chan Message
}

func NewLocalPartition(filepath string) (Partition, error) {
	outputFile, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	lp := new(localPartition)
	lp.dataFile = outputFile

	return lp, err
}

func (l localPartition) Close() {
	l.dataFile.Close()
}

func (l *localPartition) Write(payload string) error {
	thisMsgId := atomic.AddUint64(&l.currentMsgId, 1)
	m := Message{
		Id:      thisMsgId,
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
	l.listenerMutex.Lock()
	defer l.listenerMutex.Unlock()

	l.listeners = append(l.listeners, listener)
}

func informListener(listener chan Message, m Message) {
	listener <- m
}
