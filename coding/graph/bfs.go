package graph

// BFS traverses the graph breadth-first from start and returns visit order.
func BFS[T comparable](graph *Graph[T], start T) []T {
	if graph == nil || !graph.HasNode(start) {
		return nil
	}

	visited := map[T]bool{start: true}
	queue := []T{start}
	order := make([]T, 0)

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)

		for _, neighbor := range graph.Neighbors(node) {
			if visited[neighbor] {
				continue
			}

			visited[neighbor] = true
			queue = append(queue, neighbor)
		}
	}

	return order
}

func Bfs(graph *Graph[int], start int) []int {
	return BFS(graph, start)
}
