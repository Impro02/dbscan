package main

import (
	"fmt"
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
	Vect []float64
	Visited bool
	ClusterId int
}


func dbscan(points []*Point, epsilon float64, minPoints int) []int {
		clusterID := 0
	for _, point := range points {
		if point.Visited {
			continue
		}
		point.Visited = true

		neighbors := regionQuery(points, point, epsilon)
		if len(neighbors) < minPoints -1 {
			point.ClusterId = -1
		} else {
			clusterID++
			expandCluster(points, point, neighbors, clusterID, epsilon, minPoints)
		}
	}

	clusterIds := make([]int, 0, len(points))
    for _, point := range points {
        clusterIds = append(clusterIds, point.ClusterId)
    }
    return clusterIds
}

func regionQuery(points []*Point, point *Point, epsilon float64) []*Point {
	neighbors := make([]*Point, 0)

	for _, neighbor := range points {
		if point != neighbor && distance(*point, *neighbor) <= epsilon {
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}

func expandCluster(points []*Point, point *Point,  neighbors []*Point, clusterID int, epsilon float64, minPts int) {
	point.ClusterId = clusterID

	for i := 0; i < len(neighbors); i++ {
		if !neighbors[i].Visited {
			neighbors[i].Visited = true

			neighborNeighbors := regionQuery(points, neighbors[i], epsilon)

			if len(neighborNeighbors) >= minPts -1 {
				neighbors = union(neighbors, neighborNeighbors)
			}
		}

		if neighbors[i].ClusterId == 0 || neighbors[i].ClusterId == -1 {
			neighbors[i].ClusterId = clusterID
		}
	}
}

func union(a, b []*Point) []*Point {
    m := make(map[*Point]bool)

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

func distance(p1, p2 Point) float64 {
	var ret float64
	for i := 0; i < len(p1.Vect); i++ {
		tmp := p1.Vect[i] - p2.Vect[i]
		ret += tmp * tmp
	}
	return math.Sqrt(ret)
}

func main() {
	points := []*Point{
		{Vect: []float64{6, 6}},
        {Vect: []float64{1, 1}},
        {Vect: []float64{2, 2}},
		{Vect: []float64{10, 10}},
		{Vect: []float64{32, 32}},
		{Vect: []float64{33, 33}},
		{Vect: []float64{45, 45}},
		{Vect: []float64{31, 31}},
		{Vect: []float64{100, 100}},
        {Vect: []float64{3, 3}},
    }

    eps := 5.0
    minPts := 3

    clusters := dbscan(points, eps, minPts)

    for i := range clusters {
        fmt.Printf("Cluster ID: %d", i)
    }
}