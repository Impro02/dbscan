package main

import (
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
		{Vec: []float64{21, 21}},
		{Vec: []float64{3, 3}},
		{Vec: []float64{22, 22}},
		{Vec: []float64{40, 40}},
		{Vec: []float64{41, 41}},
		{Vec: []float64{20, 20}},
		{Vec: []float64{42, 42}},
		{Vec: []float64{100, 100}},
	}

	labels, clusters := dbscanGo(points, "brute", 5.0, 3)

	expectedClusters := 3
	expectedLabels := []int{1, 1, 1, -1, 2, 3, 1, 3, 2, 2, 3, 2, -1}

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
		{Vec: []float64{21, 21}},
		{Vec: []float64{3, 3}},
		{Vec: []float64{22, 22}},
		{Vec: []float64{40, 40}},
		{Vec: []float64{41, 41}},
		{Vec: []float64{20, 20}},
		{Vec: []float64{42, 42}},
		{Vec: []float64{100, 100}},
	}

	labels, clusters := dbscanGo(points, "brute", 0.5, 3)

	expectedClusters := 3
	expectedLabels := []int{1, 1, 1, -1, 2, 3, 1, 3, 2, 2, 3, 2, -1}

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
		{Vec: []float64{21, 21}},
		{Vec: []float64{3, 3}},
		{Vec: []float64{22, 22}},
		{Vec: []float64{40, 40}},
		{Vec: []float64{41, 41}},
		{Vec: []float64{20, 20}},
		{Vec: []float64{42, 42}},
		{Vec: []float64{100, 100}},
	}

	labels, clusters := dbscanGo(points, "kd_tree", 5.0, 3)

	expectedClusters := 3
	expectedLabels := []int{1, 1, 1, -1, 2, 3, 1, 3, 2, 2, 3, 2, -1}

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
		{Vec: []float64{21, 21}},
		{Vec: []float64{3, 3}},
		{Vec: []float64{22, 22}},
		{Vec: []float64{40, 40}},
		{Vec: []float64{41, 41}},
		{Vec: []float64{20, 20}},
		{Vec: []float64{42, 42}},
		{Vec: []float64{100, 100}},
	}

	labels, clusters := dbscanGo(points, "kd_tree", 0.5, 3)

	expectedClusters := 3
	expectedLabels := []int{1, 1, 1, -1, 2, 3, 1, 3, 2, 2, 3, 2, -1}

	assert.Equal(t, expectedClusters, clusters)
	assert.Equal(t, expectedLabels, labels)
}
