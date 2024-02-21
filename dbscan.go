package main

import (
	"C"
	"fmt"
	"math"
	"strconv"
)

type Point struct {
	X float64
	Y float64
}

//export DBSCAN
func DBSCAN(documentPtr *C.char) *C.char {
	documentString := C.GoString(documentPtr)

	// Parse JSON
	var data struct {
		epsilon float64  `json:"epsilon"`
		minPoints int      `json:"minPoints"`
		points []Point  `json:"points"`
	}

	err := json.Unmarshal([]byte(documentString), &data)
	if err != nil{
		log.Fatal("Error parsing JSON:", err)
		return
	}
	
	visited := make(map[Point]bool)
	clusters := [][]Point{}

	for _, point := range points {
		if visited[point] {
			continue
		}
		visited[point] = true

		neighbors := regionQuery(points, point, epsilon)
		if len(neighbors) < minPoints {
			continue
		}

		cluster := []Point{point}
		expandCluster(points, point, neighbors, epsilon, minPoints, visited, &cluster)
		clusters = append(clusters, cluster)
	}

	clustersJSON, err := json.Marshal(clusters)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil
	}

	return C.CString(string(clustersJSON))
}

func regionQuery(points []Point, point Point, epsilon float64) []Point {
	neighbors := []Point{}

	for _, p := range points {
		if distance(point, p) <= epsilon {
			neighbors = append(neighbors, p)
		}
	}

	return neighbors
}

func expandCluster(points []Point, point Point, neighbors []Point, epsilon float64, minPoints int, visited map[Point]bool, cluster *[]Point) {
	for _, neighbor := range neighbors {
		if !visited[neighbor] {
			visited[neighbor] = true

			neighborNeighbors := regionQuery(points, neighbor, epsilon)
			if len(neighborNeighbors) >= minPoints {
				neighbors = append(neighbors, neighborNeighbors...)
			}
		}

		if !isInCluster(*cluster, neighbor) {
			*cluster = append(*cluster, neighbor)
		}
	}
}

func distance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p2.X-p1.X, 2) + math.Pow(p2.Y-p1.Y, 2))
}

func isInCluster(cluster []Point, point Point) bool {
	for _, p := range cluster {
		if p == point {
			return true
		}
	}
	return false
}
