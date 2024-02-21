# DBSCAN using golang c-shared library in GO

This package implements the Density-Based Spatial Clustering of Applications with Noise (DBSCAN) algorithm using golang c-shared library in GO.

## Usage

First, import the shared library using ctypes in Python:

```
import ctypes
import json

library = ctypes.cdll.LoadLibrary("./dbscan.so")
dbscan = library.dbscan
dbscan.argtypes = [ctypes.c_char_p]
dbscan.restype = ctypes.c_void_p
```

The dbscan function takes a JSON string as input and returns a pointer to a JSON string. The input JSON should have the following structure:

```
input = {
    "epsilon": 5.0,
    "min_points": 3,
    "points": [
        {"X": 1, "Y": 1},
        {"X": 2, "Y": 2},
        {"X": 3, "Y": 3},
        {"X": 10, "Y": 10},
        {"X": 20, "Y": 20},
        {"X": 21, "Y": 21},
        {"X": 22, "Y": 22},
        {"X": 100, "Y": 100},
    ]
}
```

Here, epsilon is the maximum distance between two samples for them to be considered as in the same neighborhood, min_points is the number of samples in a neighborhood for a point to be considered as a core point, and points is the list of points to be clustered.

To call the dbscan function and parse its output, use the following code:

```
dbscan_output = dbscan(json.dumps(input).encode("utf-8"))
dbscan_output_bytes = ctypes.string_at(dbscan_output)
dbscan_result = json.loads(dbscan_output_bytes)
```

The output JSON will have the following structure:

```
{
    "Noise": [
        {"X": 10, "Y": 10},
        {"X": 100, "Y": 100},
    ],
    "Clusters": [
        [
            {"X": 1, "Y": 1},
            {"X": 2, "Y": 2},
            {"X": 3, "Y": 3},
        ],
        [
            {"X": 20, "Y": 20},
            {"X": 21, "Y": 21},
            {"X": 22, "Y": 22},
        ],
    ],
}
```

Here, Noise is the list of noise points, and Clusters is the list of clusters, where each cluster is a list of points.

