package components

// Message the container for persisted data. The ID field is considered
// unique PER partition.
type Message struct {
	ID      uint64 `json:"id"`
	Payload string `json:"payload"`
}
