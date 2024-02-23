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
    "algorithm" "kd_tree",
    "epsilon": 5.0,
    "min_points": 3,
    "points": [
        {"Vec": (1, 1)},
        {"Vec": (2, 2)},
        {"Vec": (3, 3)},
        {"Vec": (10, 10)},
        {"Vec": (20, 20)},
        {"Vec": (21, 21)},
        {"Vec": (22, 22)},
        {"Vec": (100, 100)},
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
    "labels": [1, 1, 1, -1, 2, 2, 2, -1],
    "clusters": 2,
}
```
The labels are returned in the same ordering as the given points. Noise points are labeled as -1 and are not counted as a cluster.

