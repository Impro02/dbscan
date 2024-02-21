import ctypes
import json
from unittest import TestCase

library = ctypes.cdll.LoadLibrary("./dbscan.so")


class TestDbscan(TestCase):

    @classmethod
    def setUpClass(cls) -> None:
        cls.dbscan = library.dbscan
        cls.dbscan.argtypes = [ctypes.c_char_p]
        cls.dbscan.restype = ctypes.c_void_p

        cls.points = [
            {
                "X": 1,
                "Y": 1,
            },
            {
                "X": 2,
                "Y": 2,
            },
            {
                "X": 3,
                "Y": 3,
            },
            {
                "X": 10,
                "Y": 10,
            },
            {
                "X": 20,
                "Y": 20,
            },
            {
                "X": 21,
                "Y": 21,
            },
            {
                "X": 22,
                "Y": 22,
            },
            {
                "X": 100,
                "Y": 100,
            },
        ]

    def test_dbscan(self):
        # GIVEN
        input = {
            "epsilon": 5.0,
            "min_points": 3,
            "points": self.points,
        }

        # WHEN
        dbscan_output = self.dbscan(json.dumps(input).encode("utf-8"))
        dbscan_output_bytes = ctypes.string_at(dbscan_output)
        dbscan = json.loads(dbscan_output_bytes)

        # THEN
        self.assertEqual(
            dbscan,
            {
                "Noise": [
                    {
                        "X": 10,
                        "Y": 10,
                    },
                    {
                        "X": 100,
                        "Y": 100,
                    },
                ],
                "Clusters": [
                    [
                        {
                            "X": 1,
                            "Y": 1,
                        },
                        {
                            "X": 2,
                            "Y": 2,
                        },
                        {
                            "X": 3,
                            "Y": 3,
                        },
                    ],
                    [
                        {
                            "X": 20,
                            "Y": 20,
                        },
                        {
                            "X": 21,
                            "Y": 21,
                        },
                        {
                            "X": 22,
                            "Y": 22,
                        },
                    ],
                ],
            },
        )

    def test_dbscan_only_noise(self):
        # GIVEN
        input = {
            "epsilon": 0.5,
            "min_points": 3,
            "points": self.points,
        }

        # WHEN
        dbscan_output = self.dbscan(json.dumps(input).encode("utf-8"))
        dbscan_output_bytes = ctypes.string_at(dbscan_output)
        dbscan = json.loads(dbscan_output_bytes)

        # THEN
        self.assertEqual(
            dbscan,
            {
                "Noise": [
                    {
                        "X": 1,
                        "Y": 1,
                    },
                    {
                        "X": 2,
                        "Y": 2,
                    },
                    {
                        "X": 3,
                        "Y": 3,
                    },
                    {
                        "X": 10,
                        "Y": 10,
                    },
                    {
                        "X": 20,
                        "Y": 20,
                    },
                    {
                        "X": 21,
                        "Y": 21,
                    },
                    {
                        "X": 22,
                        "Y": 22,
                    },
                    {
                        "X": 100,
                        "Y": 100,
                    },
                ],
                "Clusters": [],
            },
        )
