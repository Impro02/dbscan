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
            {"Vec": (6, 6)},
            {"Vec": (1, 1)},
            {"Vec": (2, 2)},
            {"Vec": (10, 10)},
            {"Vec": (43, 43)},
            {"Vec": (21, 21)},
            {"Vec": (3, 3)},
            {"Vec": (22, 22)},
            {"Vec": (40, 40)},
            {"Vec": (41, 41)},
            {"Vec": (20, 20)},
            {"Vec": (42, 42)},
            {"Vec": (100, 100)},
        ]

    def test_dbscan_brute(self):
        # GIVEN
        input = {
            "algorithm": "brute",
            "epsilon": 5.0,
            "min_points": 3,
            "points": self.points,
        }

        expected_output = {
            "labels": [1, 1, 1, -1, 2, 3, 1, 3, 2, 2, 3, 2, -1],
            "clusters": 3,
        }

        # WHEN
        dbscan_output = self.dbscan(json.dumps(input).encode("utf-8"))
        dbscan_output_bytes = ctypes.string_at(dbscan_output)
        dbscan = json.loads(dbscan_output_bytes)

        # THEN
        self.assertEqual(
            expected_output,
            dbscan,
        )

    def test_dbscan_brute_only_noise(self):
        # GIVEN
        input = {
            "algorithm": "brute",
            "epsilon": 0.5,
            "min_points": 3,
            "points": self.points,
        }

        expected_output = {
            "labels": [-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1],
            "clusters": 0,
        }

        # WHEN
        dbscan_output = self.dbscan(json.dumps(input).encode("utf-8"))
        dbscan_output_bytes = ctypes.string_at(dbscan_output)
        dbscan = json.loads(dbscan_output_bytes)

        # THEN
        self.assertEqual(
            expected_output,
            dbscan,
        )

    def test_dbscan_kd_tree(self):
        # GIVEN
        input = {
            "algorithm": "kd_tree",
            "epsilon": 5.0,
            "min_points": 3,
            "points": self.points,
        }

        expected_output = {
            "labels": [1, 1, 1, -1, 2, 3, 1, 3, 2, 2, 3, 2, -1],
            "clusters": 3,
        }

        # WHEN
        dbscan_output = self.dbscan(json.dumps(input).encode("utf-8"))
        dbscan_output_bytes = ctypes.string_at(dbscan_output)
        dbscan = json.loads(dbscan_output_bytes)

        # THEN
        self.assertEqual(
            expected_output,
            dbscan,
        )

    def test_dbscan_kd_tree_only_noise(self):
        # GIVEN
        input = {
            "algorithm": "kd_tree",
            "epsilon": 0.5,
            "min_points": 3,
            "points": self.points,
        }

        expected_output = {
            "labels": [-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1],
            "clusters": 0,
        }

        # WHEN
        dbscan_output = self.dbscan(json.dumps(input).encode("utf-8"))
        dbscan_output_bytes = ctypes.string_at(dbscan_output)
        dbscan = json.loads(dbscan_output_bytes)

        # THEN
        self.assertEqual(
            expected_output,
            dbscan,
        )
