package main

import "C"
import (
	"encoding/json"
	"log"

	"github.com/Impro02/kdtree"
)

type Input struct {
	Algorithm string            `json:"algorithm"`
	Epsilon   float64           `json:"epsilon"`
	MinPoints int               `json:"min_points"`
	LeafSize  int               `json:"leaf_size"`
	Points    []*EuclideanPoint `json:"points"`
	Workers   int               `json:"workers"`
}

type Output struct {
	Labels   []int `json:"labels"`
	Clusters int   `json:"clusters"`
}

type PointAndNeighbors struct {
	Point     *EuclideanPoint
	Neighbors []*EuclideanPoint
}

type EuclideanPoint struct {
	kdtree.Point
	Vec       []float64
	Visited   bool
	ClusterId int
}

func (p *EuclideanPoint) Vector() []float64 {
	return p.Vec
}

func (p *EuclideanPoint) Dim() int {
	return len(p.Vec)
}

func (p *EuclideanPoint) GetValue(dim int) float64 {
	return p.Vec[dim]
}

func (p *EuclideanPoint) Distance(other kdtree.Point) float64 {
	var ret float64
	for i := 0; i < p.Dim(); i++ {
		tmp := p.GetValue(i) - other.GetValue(i)
		ret += tmp * tmp
	}
	return ret
}

//export dbscan
func dbscan(documentPtr *C.char) *C.char {
	documentString := C.GoString(documentPtr)

	var input Input
	err := json.Unmarshal([]byte(documentString), &input)
	if err != nil {
		log.Fatal("error parsing JSON")
	}

	labels, clusterID := dbscanGo(input.Points, input.Algorithm, input.Epsilon, input.MinPoints, input.LeafSize, input.Workers)

	output := &Output{
		Labels:   labels,
		Clusters: clusterID,
	}

	outputBytes, err := json.Marshal(output)

	if err != nil {
		log.Fatal("error marshaling JSON")
	}

	outputStr := string(outputBytes)

	return C.CString(outputStr)
}

func dbscanGo(points []*EuclideanPoint, algorithm string, epsilon float64, minPoints int, leafSize int, workers int) ([]int, int) {
	var kdTree *kdtree.Node
	if algorithm == "kd_tree" {
		kdPoints := make([]kdtree.Point, 0)
		for _, p := range points {
			kdPoints = append(kdPoints, p)
		}

		kdTree = kdtree.BuildKDTree(kdPoints, 0, leafSize)
	}

	// Create a channel to send points to workers.
	pointsChan := make(chan *EuclideanPoint, len(points))

	// Create a channel to receive neighbors from workers.
	neighborsChan := make(chan PointAndNeighbors, len(points))

	// Start a number of workers.
	for i := 0; i < workers; i++ {
		go func() {
			for point := range pointsChan {
				neighbors := findNeighbors(kdTree, points, point, epsilon, algorithm)
				neighborsChan <- PointAndNeighbors{Point: point, Neighbors: neighbors}
			}
		}()
	}

	// Send all points to the workers.
	for _, point := range points {
		pointsChan <- point
	}
	close(pointsChan)

	// Receive all computedNeighbors from the workers.
	computedNeighbors := make(map[*EuclideanPoint][]*EuclideanPoint, len(points))
	for range points {
		pn := <-neighborsChan
		computedNeighbors[pn.Point] = pn.Neighbors
	}

	clusterID := 0
	for _, point := range points {
		if point.Visited {
			continue
		}
		point.Visited = true

		neighbors := computedNeighbors[point]

		if len(neighbors) < minPoints-1 {
			point.ClusterId = -1
		} else {
			clusterID++
			expandCluster(kdTree, points, point, computedNeighbors, clusterID, epsilon, minPoints, algorithm)
		}
	}

	clusterIDs := make([]int, 0, len(points))
	for _, point := range points {
		clusterIDs = append(clusterIDs, point.ClusterId)
	}
	return clusterIDs, clusterID
}

func findNeighbors(kdTree *kdtree.Node, points []*EuclideanPoint, point *EuclideanPoint, epsilon float64, algorithm string) []*EuclideanPoint {
	var neighbors = []*EuclideanPoint{}
	switch algorithm {
	case "kd_tree":
		neighbors = regionQueryKDTree(kdTree, point, epsilon)
	case "brute":
		neighbors = regionQueryBruteForce(points, point, epsilon)
	default:
		log.Fatal("invalid algorithm")
	}

	return neighbors
}

func regionQueryBruteForce(points []*EuclideanPoint, point *EuclideanPoint, epsilon float64) []*EuclideanPoint {
	neighbors := make([]*EuclideanPoint, 0)

	for _, neighbor := range points {
		if point != neighbor && point.Distance(neighbor) <= epsilon {
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}

func regionQueryKDTree(kdTree *kdtree.Node, point *EuclideanPoint, radius float64) []*EuclideanPoint {
	kdPoints := kdTree.SearchInRadius(point, radius)

	neighbors := make([]*EuclideanPoint, len(kdPoints))
	for i, p := range kdPoints {
		neighbors[i] = p.(*EuclideanPoint)
	}

	return neighbors
}

func expandCluster(kdTree *kdtree.Node, points []*EuclideanPoint, point *EuclideanPoint, computedNeighbors map[*EuclideanPoint][]*EuclideanPoint, clusterID int, epsilon float64, minPts int, algorithm string) {
	point.ClusterId = clusterID

	neighbors := computedNeighbors[point]

	queue := make([]*EuclideanPoint, len(neighbors))
	copy(queue, neighbors)

	for len(queue) > 0 {
		neighbor := queue[0]
		queue = queue[1:]

		if !neighbor.Visited {
			neighbor.Visited = true

			neighborNeighbors := computedNeighbors[neighbor]

			if len(neighborNeighbors) >= minPts-1 {
				queue = append(queue, neighborNeighbors...)
			}
		}

		if neighbor.ClusterId == 0 || neighbor.ClusterId == -1 {
			neighbor.ClusterId = clusterID
		}
	}
}

func union(a, b []*EuclideanPoint) []*EuclideanPoint {
	m := make(map[*EuclideanPoint]bool)

	// First pass to fill the map with unique elements from a and add them to result
	result := make([]*EuclideanPoint, 0, len(a)+len(b))
	for _, item := range a {
		if !m[item] {
			m[item] = true
			result = append(result, item)
		}
	}

	// Second pass to add unique elements from b to result
	for _, item := range b {
		if !m[item] {
			m[item] = true
			result = append(result, item)
		}
	}

	return result
}

func main() {}
