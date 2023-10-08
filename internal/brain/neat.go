// Package brain contains the logic for the neural network.
package brain

import (
	"fmt"
	stdrand "math/rand"
	"slices"
	"strings"
	"sync/atomic"
)

var lastInnovation atomic.Int32

// Node represents a node in a neural network.
type Node struct {
	incomming []*Connection
	NodeType  NodeType
	Bias      float64
	tempVal   float64
	ID        int
	tempSet   bool
}

func (n *Node) String() string {
	return fmt.Sprintf("[%s:%d] %+.3f <%+.3f:%t>", n.NodeType, n.ID, n.Bias, n.tempVal, n.tempSet)
}

// hasConnectionFrom returns true if the given node is connected directly or
// indirectly to this node.
func (n *Node) hasConnectionFrom(node *Node) bool {
	if node == nil {
		return false
	}
	if n.ID == node.ID {
		return true
	}
	for i := range n.incomming {
		if n.incomming[i].outNode.hasConnectionFrom(node) {
			return true
		}
	}
	return false
}

// NodeType represents the type of a node.
type NodeType int

const (
	// InputNode is a node that receives input.
	InputNode NodeType = iota
	// HiddenNode is a node that is neither an input nor an output.
	HiddenNode
	// OutputNode is a node that produces output.
	OutputNode
)

func (n NodeType) String() string {
	switch n {
	case InputNode:
		return "InputNode"
	case HiddenNode:
		return "HiddenNode"
	case OutputNode:
		return "OutputNode"
	}
	return "Unknown"
}

// Connection represents a connection between two nodes in a neural network.
type Connection struct {
	inNode     *Node
	outNode    *Node
	weight     float64
	enabled    bool
	innovation int32
}

// String returns the string representation of this connection.
func (c *Connection) String() string {
	connectivity := fmt.Sprintf("{%+.3f}", c.weight)
	if !c.enabled {
		connectivity = " / "
	}
	return fmt.Sprintf("(%s)--[%s]-->(%s)", c.inNode, connectivity, c.outNode)
}

// NEAT implements the NEAT genetic algorithm for evolving neural networks.
type NEAT struct {
	rand         *stdrand.Rand
	nodes        []*Node
	inputs       int
	outputs      int
	mutationRate int
}

// NewNEAT creates a new NEAT with the given number of input and output nodes.
// The result is a fully connected network with randomly set biases and
// weights.
func NewNEAT(inputs, outputs, mutationRate int, rand *stdrand.Rand) *NEAT {
	inputNodes := make([]*Node, inputs, inputs+outputs)
	outputNodes := make([]*Node, outputs)
	for i := range inputNodes {
		inputNodes[i] = &Node{
			NodeType: InputNode,
			ID:       i + 1,
		}

		for j := range outputNodes {
			if outputNodes[j] == nil {
				// This is the first time we are creating this node.
				outputNodes[j] = &Node{
					NodeType: OutputNode,
					ID:       len(inputNodes) + j + 1,
				}
			}
			connection := &Connection{
				inNode:     inputNodes[i],
				outNode:    outputNodes[j],
				weight:     rand.Float64() - 0.5, // [-0.5, 0.5)
				enabled:    true,
				innovation: lastInnovation.Add(1),
			}
			outputNodes[j].incomming = append(outputNodes[j].incomming, connection)
		}
	}

	return &NEAT{
		nodes:        append(inputNodes, outputNodes...),
		inputs:       inputs,
		outputs:      outputs,
		mutationRate: mutationRate,
		rand:         rand,
	}
}

func (n *NEAT) String() string {
	str := strings.Builder{}
	for i := range n.nodes {
		for j := range n.nodes[i].incomming {
			incomming := n.nodes[i].incomming[j]
			str.WriteString(incomming.String())
			str.WriteString("\n")
		}
	}
	return str.String()
}

// Predict returns the output of the neural network for the given input.
func (n *NEAT) Predict(input []float64) ([]float64, error) {
	if len(input) != n.inputs {
		return nil, fmt.Errorf("wrong input neurons, want %d got %d", n.inputs, len(input))
	}
	for i := range n.nodes {
		n.nodes[i].tempVal = 0
		n.nodes[i].tempSet = false
	}
	inputNodes := filter(n.nodes, func(n *Node) bool { return n.NodeType == InputNode })
	for i := range inputNodes {
		inputNodes[i].Bias = input[i]
		inputNodes[i].tempVal = input[i]
		inputNodes[i].tempSet = true
	}

	ret := make([]float64, n.outputs)
	lastNodes := filter(n.nodes, func(n *Node) bool { return n.NodeType == OutputNode })
	for i, node := range lastNodes {
		ret[i] = calculate(node)
	}
	return ret, nil
}

func filter(nodes []*Node, fn func(*Node) bool) []*Node {
	ret := make([]*Node, 0, len(nodes))
	for i := range nodes {
		node := nodes[i]
		if fn(node) {
			ret = append(ret, node)
		}
	}
	return ret
}

func filterPick(nodes []*Node, fn func(*Node) bool) *Node {
	nodes = filter(nodes, fn)
	index := stdrand.Intn(len(nodes))
	return nodes[index]
}

func calculate(node *Node) float64 {
	if node.tempSet || len(node.incomming) == 0 {
		return node.tempVal
	}
	val := node.Bias
	for i := range node.incomming {
		c := node.incomming[i]
		val += calculate(c.inNode) * c.weight
	}
	node.tempVal = val
	node.tempSet = true
	return val
}

// Clone returns a clone of the network.
func (n *NEAT) Clone() *NEAT {
	node := &NEAT{
		nodes:   make([]*Node, 0, len(n.nodes)),
		inputs:  n.inputs,
		outputs: n.outputs,
	}
	// we need to keep a cache of nodes in case there was a multiple connection
	// to the same node.
	nodeMap := make(map[int]*Node, len(n.nodes))

	var cloneNode func(node *Node) *Node
	cloneNode = func(node *Node) *Node {
		n := &Node{
			NodeType:  node.NodeType,
			Bias:      node.Bias,
			ID:        node.ID,
			incomming: make([]*Connection, len(node.incomming)),
		}
		for i := range node.incomming {
			tmp, ok := nodeMap[node.incomming[i].inNode.ID]
			if !ok {
				tmp = cloneNode(node.incomming[i].inNode)
			}
			n.incomming[i] = &Connection{
				inNode:     tmp,
				outNode:    n,
				weight:     node.incomming[i].weight,
				enabled:    node.incomming[i].enabled,
				innovation: node.incomming[i].innovation,
			}
			nodeMap[node.incomming[i].inNode.ID] = tmp
		}
		return n
	}

	for i := range n.nodes {
		node.nodes[i] = cloneNode(n.nodes[i])
	}
	return node
}

// Mutate mutates the NEAT. There is a 10% chance of each mutation. It might
// add a new node, delete a node, split a connection, add a new connection or
// delete a connection.
func (n *NEAT) Mutate() *NEAT {
	newNode := n.rand.Intn(100) < n.mutationRate
	deleteNode := n.rand.Intn(100) < n.mutationRate
	splitConnection := n.rand.Intn(100) < n.mutationRate
	newConnection := n.rand.Intn(100) < n.mutationRate
	deleteConnection := n.rand.Intn(100) < n.mutationRate
	toggleConnection := n.rand.Intn(100) < n.mutationRate
	changeBias := n.rand.Intn(100) < n.mutationRate
	changeWeight := n.rand.Intn(100) < n.mutationRate
	switch {
	case newNode:
		n.addRandomNode()
	case deleteNode:
		n.deleteRandomNode()
	case splitConnection:
		n.splitRandomConnection()
	case newConnection:
		n.addRandomConnection()
	case deleteConnection:
		n.deleteRandomConnection()
	case toggleConnection:
		n.toggleRandomConnection()
	case changeBias:
		n.changeRandomBias()
	case changeWeight:
		n.changeRandomWeight()
	}
	return n
}

// findNonCircularNodes finds two nodes that are not connected to each other.
func (n *NEAT) findNonCircularNodes() (inNode, outNode *Node) {
	inputNodes := filter(n.nodes, func(n *Node) bool {
		return n.NodeType == InputNode || n.NodeType == HiddenNode
	})
	outputNodes := filter(n.nodes, func(n *Node) bool {
		return n.NodeType == OutputNode || n.NodeType == HiddenNode
	})
	for {
		in := inputNodes[n.rand.Intn(len(inputNodes))]
		out := outputNodes[n.rand.Intn(len(outputNodes))]
		if !in.hasConnectionFrom(out) {
			return in, out
		}
	}
}

// addRandomNode adds a new random node to the network. It connects the node to
// random nodes. The incomming connection can be any node besides an output
// node, and the outgoing connection can be any node besides an input node. It
// prevents cycles in the network.
func (n *NEAT) addRandomNode() {
	node := &Node{
		NodeType: HiddenNode,
		Bias:     n.rand.Float64() - 0.5, // [-0.5, 0.5)
		ID:       len(n.nodes),
	}

	inNode, outNode := n.findNonCircularNodes()
	connection := &Connection{
		inNode:     inNode,
		outNode:    node,
		weight:     n.rand.Float64() - 0.5, // [-0.5, 0.5)
		enabled:    true,
		innovation: lastInnovation.Add(1),
	}
	inNode.incomming = append(inNode.incomming, connection)

	connection = &Connection{
		inNode:     node,
		outNode:    outNode,
		weight:     n.rand.Float64() - 0.5, // [-0.5, 0.5)
		enabled:    true,
		innovation: lastInnovation.Add(1),
	}
	node.incomming = append(node.incomming, connection)

	n.nodes = append(n.nodes, node)
}

// deleteRandomNode randomly deletes a node from the network. It deletes all
// the connections to this node.
func (n *NEAT) deleteRandomNode() {
	// we can't delete input or output nodes.
	node := filterPick(n.nodes, func(n *Node) bool {
		return n.NodeType == HiddenNode
	})

	// We need to remove all connections to this node.
	for i := range n.nodes {
		n.nodes[i].incomming = slices.DeleteFunc(n.nodes[i].incomming, func(c *Connection) bool {
			return c.outNode.ID == node.ID
		})
	}
	// Now we should normalise the IDs.
	for i := range n.nodes {
		n.nodes[i].ID = i + 1
	}
}

// splitRandomConnection splits a random connection in the network by adding a
// new node in between the connection.
func (n *NEAT) splitRandomConnection() {
	// We find any nodes that are not input nodes and operate on its incomming
	// connections. We don't know if the next random node has a connection, so
	// we filter for one
	node := filterPick(n.nodes, func(n *Node) bool {
		return n.NodeType == HiddenNode && len(n.incomming) > 0
	})
	connection := node.incomming[n.rand.Intn(len(node.incomming))]

	newNode := &Node{
		NodeType: HiddenNode,
		Bias:     n.rand.Float64() - 0.5, // [-0.5, 0.5)
	}
	newConnection := &Connection{
		inNode:     connection.inNode,
		outNode:    newNode,
		weight:     n.rand.Float64() - 0.5, // [-0.5, 0.5)
		enabled:    true,
		innovation: lastInnovation.Add(1),
	}
	node.incomming = append(node.incomming, newConnection)
	connection.inNode = newNode

	// Now we should normalise the IDs.
	for i := range n.nodes {
		n.nodes[i].ID = i + 1
	}
}

// addRandomConnection adds a random connection between two random nodes in the
// network.
func (n *NEAT) addRandomConnection() {
	node1, node2 := n.findNonCircularNodes()
	connection := &Connection{
		inNode:     node1,
		outNode:    node2,
		weight:     n.rand.Float64() - 0.5, // [-0.5, 0.5)
		enabled:    true,
		innovation: lastInnovation.Add(1),
	}
	node1.incomming = append(node1.incomming, connection)
}

// deleteRandomConnection deletes a random connection from the network.
func (n *NEAT) deleteRandomConnection() {
	node := filterPick(n.nodes, func(n *Node) bool {
		return n.NodeType == HiddenNode && len(n.incomming) > 0
	})

	index := n.rand.Intn(len(node.incomming))
	node.incomming = slices.Delete(node.incomming, index, index)
}

// toggleRandomConnection toggles a random connection in the network.
func (n *NEAT) toggleRandomConnection() {
	node := filterPick(n.nodes, func(n *Node) bool {
		return n.NodeType == HiddenNode && len(n.incomming) > 0
	})

	index := n.rand.Intn(len(node.incomming))
	node.incomming[index].enabled = !node.incomming[index].enabled
}

// changeRandomBias changes the bias of a random node in the network.
func (n *NEAT) changeRandomBias() {
	node := filterPick(n.nodes, func(n *Node) bool {
		return n.NodeType == HiddenNode
	})
	c := float64(1)
	if n.rand.Intn(100) > 50 {
		c = -1
	}
	node.Bias += c * 0.001
}

// changeRandomWeight changes the weight of a random connection in the network.
func (n *NEAT) changeRandomWeight() {
	node := filterPick(n.nodes, func(n *Node) bool {
		return n.NodeType == HiddenNode && len(n.incomming) > 0
	})

	index := n.rand.Intn(len(node.incomming))
	c := float64(1)
	if n.rand.Intn(100) > 50 {
		c = -1
	}
	node.incomming[index].weight += c * 0.001
}
