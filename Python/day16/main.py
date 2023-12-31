import time
from collections import deque
from enum import Enum
from typing import Deque


class Direction(Enum):
    UP = 1
    DOWN = 2
    LEFT = 3
    RIGHT = 4


class Point:
    def __init__(self, x, y):
        self.x = x
        self.y = y

    def shift(self, direction: Direction):
        if direction == Direction.UP:
            self.y -= 1
        elif direction == Direction.DOWN:
            self.y += 1
        elif direction == Direction.LEFT:
            self.x -= 1
        elif direction == Direction.RIGHT:
            self.x += 1

    def __repr__(self):
        return f"Point({self.x}, {self.y})"

    def __str__(self):
        return f"({self.x}, {self.y})"

    def tuple(self) -> tuple[int, int]:
        return self.x, self.y



class Status:

    def __init__(self, point: Point, direction: Direction):
        self.point = point
        self.direction = direction

    def move(self):
        self.point.shift(self.direction)

    def __repr__(self):
        return f"Status({self.point.x}, {self.point.y}, {self.direction})"

    def __str__(self):
        return self.__repr__()

    def tuple(self) -> tuple[int, int, int]:
        return self.point.x, self.point.y, self.direction.value


def move_point(point: (int, int), direction: int) -> tuple[int, int, int]:
    if direction == Direction.UP.value:
        return point[0], point[1] - 1, direction
    elif direction == Direction.DOWN.value:
        return point[0], point[1] + 1, direction
    elif direction == Direction.LEFT.value:
        return point[0] - 1, point[1], direction
    elif direction == Direction.RIGHT.value:
        return point[0] + 1, point[1], direction


def load_data(filename: str) -> list[list[str]]:
    data = []
    with open(filename, 'r') as file:
        for line in file:
            directions = []
            for char in line.strip():
                directions.append(char)
            data.append(directions)
    return data


def pad_grid(grid: list[list[str]]) -> list[list[str]]:
    dim = len(grid) + 2
    new_grid = [['#' for _ in range(dim)]]
    for row in grid:
        row.insert(0, '#')
        row.append('#')
        new_grid.append(row)
    new_grid.append(['#' for _ in range(dim)])
    return new_grid


def process_splitter(status: tuple[int, int, int], queue: Deque[tuple[int, int, int]],
                     direction: tuple[Direction, Direction], new_direction: tuple[Direction, Direction]):
    if status[2] == direction[0].value or status[2] == direction[1].value:
        queue.append((status[0], status[1], new_direction[0].value))
        queue.append((status[0], status[1], new_direction[1].value))
    else:
        queue.append(status)


def process_corner(status: tuple[int, int, int], queue: Deque[tuple[int, int, int]],
                   new_direction: tuple[Direction, Direction, Direction, Direction]):
    if status[2] == Direction.UP.value:
        queue.append((status[0], status[1], new_direction[0].value))
    elif status[2] == Direction.DOWN.value:
        queue.append((status[0], status[1], new_direction[1].value))
    elif status[2] == Direction.LEFT.value:
        queue.append((status[0], status[1], new_direction[2].value))
    else:
        queue.append((status[0], status[1], new_direction[3].value))


def calculate_energized_tiles(grid: list[list[str]], start: Status) -> int:
    queue: Deque[tuple[int, int, int]] = deque([start.tuple()])
    visited: set[tuple[int, int, int]] = set()
    count: set[tuple[int, int]] = set()
    while len(queue) > 0:
        current = queue.popleft()
        if current in visited:
            continue
        visited.add(current)
        count.add((current[0], current[1]))
        current = move_point(current[:2], current[2])
        if grid[current[1]][current[0]] == '-':
            process_splitter(current, queue, (Direction.UP, Direction.DOWN), (Direction.LEFT, Direction.RIGHT))
        elif grid[current[1]][current[0]] == "|":
            process_splitter(current, queue, (Direction.LEFT, Direction.RIGHT), (Direction.UP, Direction.DOWN))
        elif grid[current[1]][current[0]] == "/":
            process_corner(current, queue, (Direction.RIGHT, Direction.LEFT, Direction.DOWN, Direction.UP))
        elif grid[current[1]][current[0]] == "\\":
            process_corner(current, queue, (Direction.LEFT, Direction.RIGHT, Direction.UP, Direction.DOWN))
        elif grid[current[1]][current[0]] == ".":
            if current[2] == Direction.DOWN.value:
                current_y = current[1]
                while grid[current_y + 1][current[0]] == ".":
                    count.add((current[0], current_y))
                    current_y += 1
                queue.append((current[0], current_y, current[2]))
            elif current[2] == Direction.UP.value:
                current_y = current[1]
                while grid[current_y - 1][current[0]] == ".":
                    count.add((current[0], current_y))
                    current_y -= 1
                queue.append((current[0], current_y, current[2]))
            elif current[2] == Direction.LEFT.value:
                current_x = current[0]
                while grid[current[1]][current_x - 1] == ".":
                    count.add((current_x, current[1]))
                    current_x -= 1
                queue.append((current_x, current[1], current[2]))
            else:  # Right
                current_x = current[0]
                while grid[current[1]][current_x + 1] == ".":
                    count.add((current_x, current[1]))
                    current_x += 1
                queue.append((current_x, current[1], current[2]))
    return len(count) - 1


def print_grid(energized_points: set[tuple[int, int]]):
    dim = 12
    grid = [['.' for _ in range(dim)] for _ in range(dim)]
    for point in energized_points:
        grid[point[1]][point[0]] = '#'
    for row in grid:
        print("".join(row))
    print()


def solution_part_1():
    grid = load_data('input.txt')
    grid = pad_grid(grid)
    start = Status(Point(0, 1), Direction.RIGHT)
    energized_tiles = calculate_energized_tiles(grid, start)
    print(energized_tiles)


def solution_part_2():
    grid = pad_grid(load_data("input.txt"))
    dim = len(grid)
    energized_tiles = []
    for i in range(dim - 2):
        energized_tiles.append(calculate_energized_tiles(
            grid, Status(Point(0, i + 1), Direction.RIGHT)
        ))
        energized_tiles.append(calculate_energized_tiles(
            grid, Status(Point(dim - 1, i + 1), Direction.LEFT)
        ))
        energized_tiles.append(calculate_energized_tiles(
            grid, Status(Point(i + 1, 0), Direction.DOWN)
        ))
        energized_tiles.append(calculate_energized_tiles(
            grid, Status(Point(i + 1, dim - 1), Direction.UP)
        ))
    print(max(energized_tiles))


def main():
    start = time.time()
    solution_part_1()
    solution_part_2()
    end = time.time()
    print(f"Time {(end - start)}s")


if __name__ == "__main__":
    main()
