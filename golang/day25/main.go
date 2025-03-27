package main

import (
	"math"
	"os"
	"strings"

	"slices"
)

// Graph represents an undirected flow network.
type Graph struct {
	N        int // Number of vertices.
	Vertices map[string]int
	Capacity [][]int // capacity[u][v] is the capacity of edge u-v.
	Adj      [][]int // adj[u] contains all neighbors v of u.
}

// NewGraph initializes a graph with n vertices.
func NewGraph(n int) *Graph {
	capacity := make([][]int, n)
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		capacity[i] = make([]int, n)
		adj[i] = []int{}
	}
	var vertices = map[string]int{}
	return &Graph{n, vertices, capacity, adj}
}

func CopyGraph(g *Graph) *Graph {
	newGraph := NewGraph(g.N)
	for i := 0; i < g.N; i++ {
		for j := 0; j < g.N; j++ {
			newGraph.Capacity[i][j] = g.Capacity[i][j]
		}
		newGraph.Adj[i] = make([]int, len(g.Adj[i]))
		copy(newGraph.Adj[i], g.Adj[i])
	}
	for k, v := range g.Vertices {
		newGraph.Vertices[k] = v
	}
	return newGraph
}

func readData(filepath string) (*Graph, error) {
	edgeList, readErr := os.ReadFile(filepath)
	if readErr != nil {
		return nil, readErr
	}
	// get edge list
	edges := make(map[string][]string)
	for _, line := range strings.Split(string(edgeList), "\n") {
		if line == "" {
			continue
		}
		lineSplit := strings.Split(strings.TrimSpace(line), ":")
		edges[lineSplit[0]] = strings.Split(strings.TrimSpace(lineSplit[1]), " ")
	}
	// add reverse edges
	for k, v := range edges {
		for _, edge := range v {
			if _, ok := edges[edge]; !ok {
				edges[edge] = []string{}
			}
			if !slices.Contains(edges[edge], k) {
				edges[edge] = append(edges[edge], k)
			}
		}
	}
	graph := NewGraph(len(edges))
	maxIndex := 0
	for edge := range edges {
		graph.Vertices[edge] = maxIndex
		maxIndex++
	}
	for edge, edgesTo := range edges {
		index := graph.Vertices[edge]
		for _, edge := range edgesTo {
			edgeIndex := graph.Vertices[edge]

			graph.Adj[index] = append(graph.Adj[index], edgeIndex)
			graph.Capacity[index][edgeIndex] = 1
		}
	}
	return graph, nil
}

// bfs performs a breadth-first search on the residual graph.
// It returns true if there is a path from s to t and fills the parent slice to reconstruct the path.
func bfs(g *Graph, s, t int, parent []int) bool {
	visited := make([]bool, g.N)
	for i := range parent {
		parent[i] = -1
	}
	queue := []int{s}
	visited[s] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range g.Adj[u] {
			if !visited[v] && g.Capacity[u][v] > 0 {
				parent[v] = u
				visited[v] = true
				queue = append(queue, v)
			}
		}
	}
	return visited[t]
}

// edmondsKarp computes the maximum flow from s to t in the given graph.
func edmondsKarp(g *Graph, s, t int) int {
	parent := make([]int, g.N)
	maxFlow := 0

	// While there is a path from s to t in the residual graph:
	for bfs(g, s, t, parent) {
		// Find the bottleneck capacity along the path.
		pathFlow := math.MaxInt64
		for v := t; v != s; v = parent[v] {
			u := parent[v]
			if g.Capacity[u][v] < pathFlow {
				pathFlow = g.Capacity[u][v]
			}
		}
		// Update the capacities in the residual graph.
		for v := t; v != s; v = parent[v] {
			u := parent[v]
			g.Capacity[u][v] -= pathFlow
			g.Capacity[v][u] += pathFlow
		}
		maxFlow += pathFlow
	}
	return maxFlow
}

// minCut performs a BFS in the residual graph starting from s
// and returns a slice indicating which vertices are reachable from s.
func minCut(g *Graph, s int) []bool {
	visited := make([]bool, g.N)
	queue := []int{s}
	visited[s] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range g.Adj[u] {
			if !visited[v] && g.Capacity[u][v] > 0 {
				visited[v] = true
				queue = append(queue, v)
			}
		}
	}
	return visited
}

func reverseEdgeMap(edgeMap map[string]int) map[int]string {
	reversedEdgeMap := make(map[int]string)
	for k, v := range edgeMap {
		reversedEdgeMap[v] = k
	}
	return reversedEdgeMap
}

func countReachable(reachable []bool) int {
	count := 0
	for _, r := range reachable {
		if r {
			count++
		}
	}
	return count
}

func findThreeCut(graph *Graph) []bool {
	for i := 1; i < graph.N; i++ {
		graphCopy := CopyGraph(graph)
		edmondsKarp(graphCopy, 0, i)

		// Compute the minimum cut:
		reachable := minCut(graphCopy, 0)
		countCut := 0
		for u := 0; u < graph.N; u++ {
			if reachable[u] {
				for _, v := range graph.Adj[u] {
					if !reachable[v] {
						// Edge u -> v is in the min cut.
						countCut++
					}
				}
			}
		}
		if countCut == 3 {
			return reachable
		}
	}
	return nil
}

func solutionPart1(graph *Graph) int {
	reachable := findThreeCut(graph)
	if reachable == nil {
		return -1
	}

	return countReachable(reachable) * (graph.N - countReachable(reachable))
}

func main() {
	graph, err := readData("input.txt")
	if err != nil {
		panic(err)
	}
	println(solutionPart1(graph))
}
