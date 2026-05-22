package graph

// Graph is a generic adjacency-list graph that supports directed and undirected edges.
type Graph[T comparable] struct {
	adjacency map[T][]T
	directed  bool
}

// NewGraph creates a new generic graph.
func NewGraph[T comparable](directed bool) *Graph[T] {
	return &Graph[T]{
		adjacency: make(map[T][]T),
		directed:  directed,
	}
}

// AddNode inserts a node if it does not already exist.
func (g *Graph[T]) AddNode(node T) {
	if g == nil {
		return
	}

	if _, exists := g.adjacency[node]; !exists {
		g.adjacency[node] = []T{}
	}
}

// AddEdge inserts an edge and creates missing nodes automatically.
func (g *Graph[T]) AddEdge(from, to T) {
	if g == nil {
		return
	}

	g.AddNode(from)
	g.AddNode(to)
	g.adjacency[from] = append(g.adjacency[from], to)

	if !g.directed {
		g.adjacency[to] = append(g.adjacency[to], from)
	}
}

// HasNode reports whether the graph contains the given node.
func (g *Graph[T]) HasNode(node T) bool {
	if g == nil {
		return false
	}

	_, exists := g.adjacency[node]
	return exists
}

// Neighbors returns a copy of the node's adjacency list.
func (g *Graph[T]) Neighbors(node T) []T {
	if g == nil {
		return nil
	}

	neighbors, exists := g.adjacency[node]
	if !exists {
		return nil
	}

	return append([]T(nil), neighbors...)
}

// BuildSampleGraph returns a classic interview-style undirected graph.
func BuildSampleGraph() *Graph[int] {
	graph := NewGraph[int](false)

	edges := [][2]int{
		{0, 1},
		{0, 2},
		{1, 3},
		{1, 4},
		{2, 5},
		{3, 6},
		{4, 6},
		{5, 6},
	}

	for _, edge := range edges {
		graph.AddEdge(edge[0], edge[1])
	}

	return graph
}
