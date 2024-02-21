package main

import "C"
import (
	"encoding/json"
	"log"
	"math"
)

type Input struct {
    Epsilon   float64 `json:"epsilon"`
    MinPoints int     `json:"min_points"`
    Points    []Point `json:"points"`
}

type Output struct {
	Labels []int  `json:"labels"`
	Clusters int  `json:"clusters"`
}
type Point struct {
	X float64  `json:"X"`
	Y float64 `json:"Y"`
}


//export dbscan
func dbscan(documentPtr *C.char) *C.char {
	documentString := C.GoString(documentPtr)

	var input Input
    err := json.Unmarshal([]byte(documentString), &input)
    if err != nil {
        log.Fatal("error parsing JSON")
    }
	
	visited := make([]bool, len(input.Points))
	labels := make([]int, len(input.Points))

	clusterID := 0
	for i, point := range input.Points {
		if visited[i] {
			continue
		}
		visited[i] = true

		neighbors := regionQuery(input.Points, point, input.Epsilon)
		if len(neighbors) < input.MinPoints {
			labels[i] = -1
		} else {
			clusterID++
			expandCluster(input.Points, visited, labels, i, neighbors, clusterID, input.Epsilon, input.MinPoints)
		}
	}

	output := &Output{
		Labels: labels, 
		Clusters: clusterID,
	}

	outputBytes, err := json.Marshal(output)

	if err != nil {
		log.Fatal("error marshaling JSON")
	}

	outputStr := string(outputBytes)

	return C.CString(outputStr)
}

func regionQuery(points []Point, point Point, epsilon float64) []int {
	neighbors := []int{}

	for i, neighbor := range points {
		if distance(point, neighbor) <= epsilon {
			neighbors = append(neighbors, i)
		}
	}

	return neighbors
}

func expandCluster(points []Point, visited []bool, labels []int, pointIndex int, neighbors []int, clusterID int, epsilon float64, minPts int) {
	labels[pointIndex] = clusterID

	for _, neighborIndex := range neighbors {
		if !visited[neighborIndex] {
			visited[neighborIndex] = true

			neighborNeighbors := regionQuery(points, points[neighborIndex], epsilon)

			if len(neighborNeighbors) >= minPts {
				neighbors = append(neighbors, neighborNeighbors...)
			}
		}

		if labels[neighborIndex] == 0 || labels[neighborIndex] == -1 {
			labels[neighborIndex] = clusterID
		}
	}
}

func distance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p2.X-p1.X, 2) + math.Pow(p2.Y-p1.Y, 2))
}

func main(){}
