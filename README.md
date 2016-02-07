# shika-go

A simple go implementation of a kafka-esque message bus.

**Run the example**: `go run example.go  junk-topic config.json`

## State of Affairs
This currently has a single node, no recovery implementation working. That is to say that it has a lot of the plumbing for it will need, most of the distributed part is pending.

Next pieces of work:
  - create the client library that will allow the clients to not be part of the same binary distro
    - proxy restful calls to a well known API
  - implement the idea of a remote partition
    - unclear on the best way to 'subscribe' to a remote partition. Probably going to 'install' a consumer on the node with the actual partition that will just proxy the messages back to this node which will pass it through the actual channel. 
  - add a central coordination node (master node) for
    - creating and distributing partitions
    - cataloging consumer offsets
    - central node for now b/c of simplicity, eventually deal with leader election scenarios and health
  - use memmapped files for data storage
    - consider how these files will be partitioned, expired and compressed
  - deal with restarting the nodes and service
    - this means that we will re-populate the message queues, and avoid too much re-sending of messages
  - health checks for nodes and central command
  - write replication
  - remove central command? figure out how to do service discovery?
  
## How it is structured
There is the concept of a 'Node', 'topic' and 'partition'. The actual client will bond up to a node which will proxy the calls to the actual partition. If a topic hasn't been created, it will be created (TODO: distribute this creation across the cluster) and given the default number of partitions. Currently, this is a hardcoded value in the configuration but that will change. 

A node will handle the creation and communication of messages, but the partition is responsible for writing the data to disk and then pushing that data to the consumer. Currently that is via a go channel, but it could be abstracted to an RPC-esque interface for non-local consumers. Right now it all runs as part of the same binary.

IDs for `Messages` are unique *to that partition*. For a topic you can define the `RoutingStrategy` such that it will direct the messages properly. Default is Round Robin from the perspective of the node.

## Using it
Creating the Node
```
  import c "github.com/squirrely/shikago/components"
  
  func main () {
	  var config c.Configuration
    // ...populate that config (json unmarshal)
    node := c.NewNode(&config)
  }
```

Subscribing to a topic
```
  msgChannel := node.Subscribe("sometopic") // bonds to a single partition
  allChannel := node.SubscribeToAll("sometopic") // bonds to all the partitions
```

Publishing a message
```
  node.Write("sometopic", "some message") // writes that string to a single partition based on the RoutingStrategy defined
```
