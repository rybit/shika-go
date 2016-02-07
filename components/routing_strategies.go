package components

import "sync/atomic"

// RoutingStrategy Defines which partition should get this piece of data
type RoutingStrategy interface {
	WhichPartition(payload string) int32
}

type roundRobinStrategy struct {
	howManyPartitions  int32
	nextPartitionIndex int32
}

// NewRoundRobinStrategy creates a strategy that will just cycle through partitions
func NewRoundRobinStrategy(partitionCount int) RoutingStrategy {
	rr := new(roundRobinStrategy)
	rr.howManyPartitions = int32(partitionCount)

	return rr
}

func (rr *roundRobinStrategy) WhichPartition(payload string) int32 {
	parCount := atomic.AddInt32(&rr.nextPartitionIndex, 1)
	return parCount % rr.howManyPartitions
}
