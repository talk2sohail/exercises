package main

import (
	"coding/graph"
	"fmt"
)

func main() {
	g := graph.BuildSampleGraph()

	visit_order := graph.Bfs(g, 0)
	fmt.Println(visit_order)
}
