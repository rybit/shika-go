package components

// Client This is responsible for contacting a node. It is also responsible for
// dealing with the writes in such a way that they are "fair"
// type Client struct {
// 	name        string
// 	strategyMap map[string]RoutingStrategy
// 	node        *Node
// 	Listener    chan Message
// }
//
// // Producer todo
// type Producer struct {
// 	name string
// }

//
// // NewClient Creates a new client that is accessing the specified node
// func NewClient(name string, node *Node) *Client {
// 	client := new(Client)
// 	client.node = node
// 	client.strategyMap = make(map[string]RoutingStrategy)
// 	return client
// }
//
// // Publish will write to the node according to the strategy specified
// func (c Client) Publish(topic, payload string) error {
// 	strat := c.getStrategyFor(topic)
// 	partIndex := strat.WhichPartition(payload)
//
// 	return c.node.writeToPartition(topic, partIndex, payload)
// }
//
// // Subscribe will defer to the node to choose which partition to bind to
// func (c *Client) Subscribe(topic string) error {
// 	c.node.subscribe(topic, c)
// 	return nil
// }
//
// // SubscribeToAll will bind to all of the partitions the node has
// func (c Client) SubscribeToAll(topic string) error {
// 	parts := c.node.getPartitionsFor(topic)
//
// 	for _, p := range parts {
// 		p.Subscribe(c.Listener)
// 	}
//
// 	return nil
// }
//
// func (c *Client) getStrategyFor(topic string) RoutingStrategy {
// 	strat, hasStrat := c.strategyMap[topic]
// 	if hasStrat {
// 		return strat
// 	}
//
// 	parts := c.node.getPartitionsFor(topic)
//
// 	// TODO - build other types or check on what they should be
// 	strat = NewRoundRobinStrategy(len(parts))
// 	c.strategyMap[topic] = strat
// 	return strat
// }
