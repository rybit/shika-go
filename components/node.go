package components

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	s "strings"

	l "github.com/squirrely/shika-go/logging"
)

var log = l.NewLogger("Node")

// Node is the primary communication endpoint for shikago. Writes and reads are done against
// a node. The will request keep a list of the local topic/paritions that manages. It is also
// responsible for creating new Topics and paritions.
// It will manage the local server that will bind ports for the local partition
type Node struct {
	config *Configuration

	// TODO always sanitize the topic name - lowercase it
	partitions  map[string][]partition
	strategyMap map[string]RoutingStrategy
	server      http.Server
	port        int
}

// NewNode creates a new node base around this configuration
func NewNode(config *Configuration) *Node {
	node := new(Node)

	node.config = config
	node.partitions = make(map[string][]partition)
	node.strategyMap = make(map[string]RoutingStrategy)

	// create a http listener, do it this way to get the running port
	addrStr := fmt.Sprintf(":%d", config.NodePort)
	listener, err := net.Listen("tcp", addrStr)
	if err != nil {
		log.Err("Failed to create listener at address %s", addrStr, err)
	}

	node.port = listener.Addr().(*net.TCPAddr).Port
	node.server = http.Server{
		Handler: node,
	}

	log.Info("Node is starting to listen on port: %d", node.port)
	go node.server.Serve(listener)

	return node
}

// Shutdown does what is on the box
func (n Node) Shutdown() {
	log.Info("Starting shutdown of %d partitions", len(n.partitions))
	for _, parts := range n.partitions {
		for _, part := range parts {
			part.Close()
		}
	}

	log.Info("Finished closing all the partitions")
}

// Subscribe will create a channel that will be published to each time a change on ONE
// partition changes
func (n *Node) Subscribe(topic string) <-chan Message {
	consumer := make(chan Message)

	parts := n.getPartitionsFor(topic)
	var part partition
	for _, p := range parts {
		part = smallestOf(part, p)
	}

	log.Debug("Created subscriber for %s:%v", topic, part)
	part.Subscribe(consumer)
	return consumer
}

// SubscribeToAll creates a channel that is registered to ALL of the different partitions
func (n *Node) SubscribeToAll(topic string) <-chan Message {
	consumer := make(chan Message)

	parts := n.getPartitionsFor(topic)
	for _, p := range parts {
		p.Subscribe(consumer)
	}

	return consumer
}

// RegisterStrategy allows a different RoutingStrategy to be specified for a given topic
func (n *Node) RegisterStrategy(topic string, strategy RoutingStrategy) {
	n.strategyMap[topic] = strategy
}

// Write write the payload to a parition in the topic based on the RoutingStrategy specified
func (n *Node) Write(topic, payload string) error {
	rs := n.getStrategyFor(topic)
	partID := rs.WhichPartition(payload)

	parts := n.getPartitionsFor(topic)
	return parts[partID].Write(payload)
}

// -------------------------------------------------------------------------------------------------
// Utilities
// -------------------------------------------------------------------------------------------------

func (n *Node) getStrategyFor(topic string) RoutingStrategy {
	rs, exists := n.strategyMap[topic]
	if exists {
		return rs
	}

	parts := n.getPartitionsFor(topic)

	rs = NewRoundRobinStrategy(len(parts))
	n.strategyMap[topic] = rs
	return rs
}

// responsible for (1) discovering if we have seen this topic before, otherwise fetching it
func (n *Node) getPartitionsFor(topic string) []partition {
	parts, hasParts := n.partitions[topic]

	if hasParts {
		return parts
	}

	// we haven't seen this before fetch it
	log.Info("The topic %s is new -- creating it locally", topic)

	// TODO fetch!!! - for now, we will just make it all local
	parts = make([]partition, n.config.DefaultPartitionCount)
	for i := range parts {
		partName := fmt.Sprintf("%s_%d.jsonl", topic, i)
		partName = s.ToLower(s.TrimSpace(s.Replace(partName, " ", "_", -1)))

		lp, err := newLocalPartition(filepath.Join(n.config.DataDirectory, partName), topic)
		if err != nil {
			log.Err("Failed to create partition %d for topic %s", i, topic, err)
		}

		log.Debug("Created local partition topic %s, %s", topic, partName)
		parts[i] = lp
	}

	n.partitions[topic] = parts
	return parts
}

func smallestOf(existing, possible partition) partition {
	if existing == nil {
		return possible
	}

	if existing.SubscribersCount() > possible.SubscribersCount() {
		return possible
	}

	return existing
}
