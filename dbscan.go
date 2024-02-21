package main

import "C"
import (
	"encoding/json"
	"log"
	"math"
)

type Point struct {
	X float64
	Y float64
}

type Input struct {
	epsilon float64
	minPoints int
	points []Point
}

type Output struct {
	Noise []Point 
	Clusters [][]Point
}

//export dbscan
func dbscan(documentPtr *C.char) *C.char {
	documentString := C.GoString(documentPtr)

	var data map[string]interface{}

	err := json.Unmarshal([]byte(documentString), &data)
	if err != nil{
		log.Fatal("Error parsing JSON:", err)
	}

	epsilon, epsilonOk := data["epsilon"].(float64)
	minPointsFloat, minPointsOk := data["min_points"].(float64)
	pointsInterface, pointsInterfaceOk := data["points"].([]interface{})

	if !epsilonOk || !minPointsOk || !pointsInterfaceOk {
		log.Fatal("Error parsing JSON: epsilon (float64), min_points (float64) and points ([]interface{}) are required!")
	}

	minPoints := int(minPointsFloat)

	points := make([]Point, len(pointsInterface))
	for i, pointInterface := range pointsInterface {
		pointMap, ok := pointInterface.(map[string]interface{})
		if !ok {
			log.Fatal("point is not a map[string]interface{}")
		}

		x, ok := pointMap["X"].(float64)
		if !ok {
			log.Fatal("X is not a float64")
		}

		y, ok := pointMap["Y"].(float64)
		if !ok {
			log.Fatal("Y is not a float64")
		}

		points[i] = Point{X: x, Y: y}
	}
	
	visited := make(map[Point]bool)
	clusters := make([][]Point, 0)
	noise := make([]Point, 0)

	for _, point := range points {
		if visited[point] {
			continue
		}
		visited[point] = true

		neighbors := regionQuery(points, point, epsilon)
		if len(neighbors) < minPoints {
			noise = append(noise, point)
		} else {
			cluster := []Point{point}
			expandCluster(points, point, neighbors, epsilon, minPoints, visited, &cluster)
			clusters = append(clusters, cluster)
		}
	}

	output := &Output{
		Noise: noise, 
		Clusters: clusters,
	}

	outputBytes, err := json.Marshal(output)
	if err != nil {
		log.Fatal("Error marshaling JSON:", err)
	}

	outputStr := string(outputBytes)

	return C.CString(outputStr)
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

func main(){}
