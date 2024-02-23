package main

import "C"
import (
	"encoding/json"
	"log"
	"math"

	"github.com/hongshibao/go-kdtree"
)

type Input struct {
	Algorithm string            `json:"algorithm"`
	Epsilon   float64           `json:"epsilon"`
	MinPoints int               `json:"min_points"`
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
	return math.Sqrt(ret)
}

func (p *EuclideanPoint) PlaneDistance(val float64, dim int) float64 {
	tmp := p.GetValue(dim) - val
	return tmp * tmp
}

//export dbscan
func dbscan(documentPtr *C.char) *C.char {
	documentString := C.GoString(documentPtr)

	var input Input
	err := json.Unmarshal([]byte(documentString), &input)
	if err != nil {
		log.Fatal("error parsing JSON")
	}

	labels, clusterID := dbscanGo(input.Points, input.Algorithm, input.Epsilon, input.MinPoints)

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

func dbscanGo(points []*EuclideanPoint, algorithm string, epsilon float64, minPoints int) ([]int, int) {
	var kdTree *kdtree.KDTree
	if algorithm == "kd_tree" {
		kdPoints := make([]kdtree.Point, 0)
		for _, p := range points {
			kdPoints = append(kdPoints, p)
		}

		kdTree = kdtree.NewKDTree(kdPoints)
	}

	clusterID := 0
	for _, point := range points {
		if point.Visited {
			continue
		}
		point.Visited = true

		var neighbors = []*EuclideanPoint{}
		switch algorithm {
		case "kd_tree":
			neighbors = regionQueryKDTree(kdTree, point, epsilon, minPoints)
		case "brute":
			neighbors = regionQueryBruteForce(points, point, epsilon)
		default:
			log.Fatal("invalid algorithm")
		}

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

func regionQueryBruteForce(points []*EuclideanPoint, point *EuclideanPoint, epsilon float64) []*EuclideanPoint {
	neighbors := make([]*EuclideanPoint, 0)

	for _, neighbor := range points {
		if point != neighbor && point.Distance(neighbor) <= epsilon {
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}

func regionQueryKDTree(kdTree *kdtree.KDTree, point *EuclideanPoint, radius float64, k int) []*EuclideanPoint {
	kdPoints := kdTree.KNN(point, k+1)
	points := make([]*EuclideanPoint, 0)
	for _, p := range kdPoints {
		points = append(points, p.(*EuclideanPoint))
	}

	neighbors := make([]*EuclideanPoint, 0)
	for _, neighbor := range points {
		if point != neighbor && point.Distance(neighbor) <= radius {
			neighbors = append(neighbors, neighbor)
		}
	}
	return neighbors
}

func expandCluster(kdTree *kdtree.KDTree, points []*EuclideanPoint, point *EuclideanPoint, neighbors []*EuclideanPoint, clusterID int, epsilon float64, minPts int, algorithm string) {
	point.ClusterId = clusterID

	for i := 0; i < len(neighbors); i++ {
		if !neighbors[i].Visited {
			neighbors[i].Visited = true

			var neighborNeighbors = []*EuclideanPoint{}
			switch algorithm {
			case "kd_tree":
				neighborNeighbors = regionQueryKDTree(kdTree, neighbors[i], epsilon, minPts)
			case "brute":
				neighborNeighbors = regionQueryBruteForce(points, neighbors[i], epsilon)
			default:
				log.Fatal("invalid algorithm")
			}

			if len(neighborNeighbors) >= minPts-1 {
				neighbors = union(neighbors, neighborNeighbors)
			}
		}

		if neighbors[i].ClusterId == 0 || neighbors[i].ClusterId == -1 {
			neighbors[i].ClusterId = clusterID
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
