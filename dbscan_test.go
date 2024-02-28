package main

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbscanGoBrute(t *testing.T) {
	points := []*EuclideanPoint{
		{Vec: []float64{6, 6}},
		{Vec: []float64{1, 1}},
		{Vec: []float64{2, 2}},
		{Vec: []float64{10, 10}},
		{Vec: []float64{43, 43}},
		{Vec: []float64{9, 9}},
		{Vec: []float64{21, 21}},
		{Vec: []float64{3, 3}},
		{Vec: []float64{22, 22}},
		{Vec: []float64{40, 40}},
		{Vec: []float64{41, 41}},
		{Vec: []float64{20, 20}},
		{Vec: []float64{42, 42}},
		{Vec: []float64{100, 100}},
	}

	labels, clusters := dbscanGo(points, "brute", math.Pow(5.0, 2), 3, 3, 1)

	expectedClusters := 3
	expectedLabels := []int{1, 1, 1, 1, 2, 1, 3, 1, 3, 2, 2, 3, 2, -1}

	assert.Equal(t, expectedClusters, clusters)
	assert.Equal(t, expectedLabels, labels)
}

func TestDbscanGoBruteOnlyNoise(t *testing.T) {
	points := []*EuclideanPoint{
		{Vec: []float64{6, 6}},
		{Vec: []float64{1, 1}},
		{Vec: []float64{2, 2}},
		{Vec: []float64{10, 10}},
		{Vec: []float64{43, 43}},
		{Vec: []float64{9, 9}},
		{Vec: []float64{21, 21}},
		{Vec: []float64{3, 3}},
		{Vec: []float64{22, 22}},
		{Vec: []float64{40, 40}},
		{Vec: []float64{41, 41}},
		{Vec: []float64{20, 20}},
		{Vec: []float64{42, 42}},
		{Vec: []float64{100, 100}},
	}

	labels, clusters := dbscanGo(points, "brute", math.Pow(0.5, 2), 3, 2, 1)

	expectedClusters := 0
	expectedLabels := []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}

	assert.Equal(t, expectedClusters, clusters)
	assert.Equal(t, expectedLabels, labels)
}

func TestDbscanGoKdTree(t *testing.T) {
	points := []*EuclideanPoint{
		{Vec: []float64{6, 6}},
		{Vec: []float64{1, 1}},
		{Vec: []float64{2, 2}},
		{Vec: []float64{10, 10}},
		{Vec: []float64{43, 43}},
		{Vec: []float64{9, 9}},
		{Vec: []float64{21, 21}},
		{Vec: []float64{3, 3}},
		{Vec: []float64{22, 22}},
		{Vec: []float64{40, 40}},
		{Vec: []float64{41, 41}},
		{Vec: []float64{20, 20}},
		{Vec: []float64{42, 42}},
		{Vec: []float64{100, 100}},
	}

	labels, clusters := dbscanGo(points, "kd_tree", math.Pow(5.0, 2), 3, 2, 1)

	expectedClusters := 3
	expectedLabels := []int{1, 1, 1, 1, 2, 1, 3, 1, 3, 2, 2, 3, 2, -1}

	assert.Equal(t, expectedClusters, clusters)
	assert.Equal(t, expectedLabels, labels)
}

func TestDbscanGoKdTreeOnlyNoise(t *testing.T) {
	points := []*EuclideanPoint{
		{Vec: []float64{6, 6}},
		{Vec: []float64{1, 1}},
		{Vec: []float64{2, 2}},
		{Vec: []float64{10, 10}},
		{Vec: []float64{43, 43}},
		{Vec: []float64{9, 9}},
		{Vec: []float64{21, 21}},
		{Vec: []float64{3, 3}},
		{Vec: []float64{22, 22}},
		{Vec: []float64{40, 40}},
		{Vec: []float64{41, 41}},
		{Vec: []float64{20, 20}},
		{Vec: []float64{42, 42}},
		{Vec: []float64{100, 100}},
	}

	labels, clusters := dbscanGo(points, "kd_tree", math.Pow(0.5, 2), 3, 1, 1)

	expectedClusters := 0
	expectedLabels := []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}

	assert.Equal(t, expectedClusters, clusters)
	assert.Equal(t, expectedLabels, labels)
}

func TestDbscanGoKdTreeLargeDataset(t *testing.T) {
	// Seed the random number generator with a constant value.
	r := rand.New(rand.NewSource(42))

	// Define the centroids.
	centroids := []EuclideanPoint{
		{Vec: []float64{0.1, 0.1}},
		{Vec: []float64{0.1, 0.9}},
		{Vec: []float64{0.5, 0.5}},
		{Vec: []float64{0.9, 0.1}},
		{Vec: []float64{0.9, 0.9}},
	}

	points := []*EuclideanPoint{}
	for i := 0; i < 100000; i++ {
		centroid := centroids[r.Intn(len(centroids))]
		x := centroid.GetValue(0) + r.NormFloat64()*0.1
		y := centroid.GetValue(1) + r.NormFloat64()*0.1
		points = append(points, &EuclideanPoint{Vec: []float64{x, y}})
	}

	_, clusters := dbscanGo(points, "kd_tree", math.Pow(0.01, 2), 3, 2, 8)

	expectedClusters := 353

	assert.Equal(t, expectedClusters, clusters)
}
