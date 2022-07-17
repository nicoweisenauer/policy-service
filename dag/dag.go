package dag

import (
	"errors"
)

// a topological ordering represents a valid sequence for tasks (nodes) which depend on each other (edges between nodes)
func TopologicalSort(graph map[string][]string) ([][]string, error) {
	sccSlice := tarjans(graph)
	isDAG := assertDAG(sccSlice, len(graph))

	if !isDAG {
		return make([][]string, 0), errors.New("Cyclic dependency detected!")
	}

	return sccSlice, nil
}

// a DAG does not contain cycles, which means that every SCC (strongly connected component) consists of a single node
func assertDAG(sccSlice [][]string, nodeCount int) bool {
	if len(sccSlice) != nodeCount {
		return false
	}

	for _, scc := range sccSlice {
		if len(scc) != 1 {
			return false
		}
	}

	return true
}

// See: https://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm
// Tarjan's creates a slice where each item is a slice of strongly connected nodes
// If a slice item contains only one node there are no cycles. A cycle on the node itself is also a connected group.
func tarjans(graph map[string][]string) [][]string {
	context := &data{
		graph: graph,
		nodes: make([]node, 0, len(graph)),
		index: make(map[string]int, len(graph)),
	}
	for node := range context.graph {
		if _, alreadyVisited := context.index[node]; !alreadyVisited {
			context.strongConnect(node)
		}
	}
	return context.output
}

// data contains the context of the algorithm
type data struct {
	graph  map[string][]string
	nodes  []node
	stack  []string
	index  map[string]int
	output [][]string
}

// node stores data for a single node in the connection process
type node struct {
	lowlink int
	stacked bool
}

// strongConnect runs recursively and outputs a grouping of strongly connected nodes
func (data *data) strongConnect(currentNode string) *node {
	index := len(data.nodes)
	data.index[currentNode] = index
	data.stack = append(data.stack, currentNode)
	data.nodes = append(data.nodes, node{lowlink: index, stacked: true})
	node := &data.nodes[index]

	for _, successorNode := range data.graph[currentNode] {
		i, seen := data.index[successorNode]
		if !seen {
			n := data.strongConnect(successorNode)
			if n.lowlink < node.lowlink {
				node.lowlink = n.lowlink
			}
		} else if data.nodes[i].stacked {
			if i < node.lowlink {
				node.lowlink = i
			}
		}
	}

	if node.lowlink == index {
		var nodes []string
		i := len(data.stack) - 1
		for {
			stronglyConnectedNode := data.stack[i]
			stackIndex := data.index[stronglyConnectedNode]
			data.nodes[stackIndex].stacked = false
			nodes = append(nodes, stronglyConnectedNode)
			if stackIndex == index {
				break
			}
			i--
		}
		data.stack = data.stack[:i]
		data.output = append(data.output, nodes)
	}

	return node
}
