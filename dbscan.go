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
}

type Output struct {
	Labels   []int `json:"labels"`
	Clusters int   `json:"clusters"`
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

	labels, clusterID := dbscanGo(input.Points, input.Algorithm, input.Epsilon, input.MinPoints, input.LeafSize)

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

func dbscanGo(points []*EuclideanPoint, algorithm string, epsilon float64, minPoints int, leafSize int) ([]int, int) {
	var kdTree *kdtree.Node
	if algorithm == "kd_tree" {
		kdPoints := make([]kdtree.Point, 0)
		for _, p := range points {
			kdPoints = append(kdPoints, p)
		}

		kdTree = kdtree.BuildKDTree(kdPoints, 0, leafSize)
	}

	clusterID := 0
	for _, point := range points {
		if point.Visited {
			continue
		}
		point.Visited = true

		neighbors := findNeighbors(kdTree, points, point, epsilon, algorithm)

		if len(neighbors) < minPoints-1 {
			point.ClusterId = -1
		} else {
			clusterID++
			expandCluster(kdTree, points, point, neighbors, clusterID, epsilon, minPoints, algorithm)
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

func expandCluster(kdTree *kdtree.Node, points []*EuclideanPoint, point *EuclideanPoint, neighbors []*EuclideanPoint, clusterID int, epsilon float64, minPts int, algorithm string) {
	point.ClusterId = clusterID

	for i := 0; i < len(neighbors); i++ {
		neighbor := neighbors[i]
		if !neighbor.Visited {
			neighbor.Visited = true

			neighborNeighbors := findNeighbors(kdTree, points, neighbor, epsilon, algorithm)

			if len(neighborNeighbors) >= minPts-1 {
				neighbors = union(neighbors, neighborNeighbors)
			}
		}

		if neighbor.ClusterId == 0 || neighbor.ClusterId == -1 {
			neighbor.ClusterId = clusterID
		}
	}
}

func union(a, b []*EuclideanPoint) []*EuclideanPoint {
	m := make(map[*EuclideanPoint]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			a = append(a, item)
		}
	}

	return a
}

func main() {}
