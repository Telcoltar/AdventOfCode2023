from typing import Callable


def load_data(filename: str) -> list[list[str]]:
    data = []
    with open(filename, 'r') as file:
        current_map = []
        for line in file:
            if line == "\n":
                data.append(current_map)
                current_map = []
                continue
            current_map.append(line.strip())
    data.append(current_map)
    return data


def get_starting_points(data: list[str], equal_fn: Callable[[str, str], bool]) -> list[int]:
    starting_points = []
    for i in range(len(data) - 1):
        if equal_fn(data[i], data[i + 1]):
            starting_points.append(i)
    return starting_points


def is_starting_point_possible(grid: list[str], starting_point: int, equal_fn: Callable[[str, str], bool]) -> bool:
    lower = starting_point
    upper = starting_point + 1
    while lower >= 0 and upper < len(grid) and equal_fn(grid[lower], grid[upper]):
        lower -= 1
        upper += 1
    return lower == -1 or upper == len(grid)


def scan_horizontal_reflection_line(grid: list[str], equal_fn: Callable[[str, str], bool]) -> list[int]:
    starting_points = get_starting_points(grid, equal_fn)
    return [point for point in starting_points
            if is_starting_point_possible(grid, point, equal_fn)]


def transpose_grid(grid: list[str]) -> list[str]:
    transposed_grid = []
    for j in range(len(grid[0])):
        transposed_row = ""
        for i in range(len(grid)):
            transposed_row += grid[i][j]
        transposed_grid.append(transposed_row)
    return transposed_grid


def is_strict_equal(a, b: str) -> bool:
    return a == b


def is_similar(a, b: str) -> bool:
    diff = 0
    for char_a, char_b in zip(a, b):
        if char_a != char_b:
            diff += 1
    return diff == 0 or diff == 1


def solution_part_1():
    grids = load_data("input.txt")
    horizontal_sum = 0
    vertical_sum = 0
    for grid in grids:
        horizontal_points = scan_horizontal_reflection_line(grid, is_strict_equal)
        if len(horizontal_points) == 1:
            horizontal_sum += horizontal_points[0] + 1
        else:
            transposed_grid = transpose_grid(grid)
            vertical_points = scan_horizontal_reflection_line(transposed_grid, is_strict_equal)
            if len(vertical_points) == 1:
                vertical_sum += vertical_points[0] + 1
            else:
                raise "Not found exactly one reflection line"
    total = vertical_sum + horizontal_sum * 100
    print(total)


def solution_part_2():
    grids = load_data("input.txt")
    horizontal_sum = 0
    vertical_sum = 0
    for grid in grids:
        horizontal_points = scan_horizontal_reflection_line(grid, is_similar)
        horizontal_strict_point = scan_horizontal_reflection_line(grid, is_strict_equal)
        horizontal_diff = list(filter(lambda point: point not in horizontal_strict_point, horizontal_points))
        if len(horizontal_diff) == 1:
            horizontal_sum += horizontal_diff[0] + 1
        else:
            transposed_grid = transpose_grid(grid)
            vertical_points = scan_horizontal_reflection_line(transposed_grid, is_similar)
            vertical_strict_points = scan_horizontal_reflection_line(transposed_grid, is_strict_equal)
            vertical_diff = list(filter(lambda point: point not in vertical_strict_points, vertical_points))
            if len(vertical_diff) == 1:
                vertical_sum += vertical_diff[0] + 1
            else:
                raise "Not found exactly one"
    total_sum = vertical_sum + horizontal_sum * 100
    print(total_sum)


def main():
    solution_part_1()
    solution_part_2()


if __name__ == "__main__":
    main()
