package components

import "time"

// Message the container for persisted data. The ID field is considered
// unique PER partition.
type Message struct {
	Timestamp time.Time `json:"timestamp"`
	ID        uint64    `json:"id"`
	Payload   string    `json:"payload"`
}
