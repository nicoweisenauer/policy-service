package dag

// Connections creates a slice where each item is a slice of strongly connected vertices.
//
// If a slice item contains only one vertex there are no loops. A loop on the
// vertex itself is also a connected group.
func Connections(graph map[string][]string) [][]string {
	g := &data{
		graph: graph,
		nodes: make([]node, 0, len(graph)),
		index: make(map[string]int, len(graph)),
	}
	for v := range g.graph {
		if _, ok := g.index[v]; !ok {
			g.strongConnect(v)
		}
	}
	return g.output
}

// data contains all common data for a single operation.
type data struct {
	graph  map[string][]string
	nodes  []node
	stack  []string
	index  map[string]int
	output [][]string
}

// node stores data for a single vertex in the connection process.
type node struct {
	lowlink int
	stacked bool
}

// strongConnect runs Tarjan's algorithm recursively and outputs a grouping of
// strongly connected vertices.
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
		var vertices []string
		i := len(data.stack) - 1
		for {
			w := data.stack[i]
			stackIndex := data.index[w]
			data.nodes[stackIndex].stacked = false
			vertices = append(vertices, w)
			if stackIndex == index {
				break
			}
			i--
		}
		data.stack = data.stack[:i]
		data.output = append(data.output, vertices)
	}

	return node
}
