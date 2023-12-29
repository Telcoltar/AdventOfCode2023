ROUND = "O".encode()[0]
SQUARE = "#".encode()[0]
EMPTY = ".".encode()[0]


def load_data(filename: str) -> list[bytearray]:
    with open(filename, encoding="UTF-8") as f:
        grid = []
        for line in f.readlines():
            grid.append(bytearray(line.strip().encode()))
        return grid


def tilt_north(grid: list[bytearray]):
    for row in range(len(grid)):
        for col in range(len(grid)):
            if grid[row][col] == ROUND:
                current_row = row
                while current_row > 0 and grid[current_row - 1][col] == EMPTY:
                    current_row -= 1
                if current_row != row:
                    grid[row][col] = EMPTY
                    grid[current_row][col] = ROUND


def calculate_load(grid: list[bytearray]) -> int:
    load = 0
    dim = len(grid)
    for row in range(dim):
        for col in range(dim):
            if grid[row][col] == ROUND:
                load += (dim - row)
    return load


def rotate_grid(grid: list[bytearray]):
    dim = len(grid)
    for row in range(dim):
        for col in range(row, dim):
            grid[row][col], grid[col][row] = grid[col][row], grid[row][col]
    for col in range(dim // 2):
        for row in range(dim):
            grid[row][col], grid[row][dim - col - 1] = grid[row][dim - col - 1], grid[row][col]


def print_grid(grid: list[bytearray]):
    for row in grid:
        print(row)
    print()


def grid_to_string(grid: list[bytearray]):
    return "".join([row.decode() for row in grid])


def solution_part_1():
    grid = load_data("input_example.txt")
    tilt_north(grid)
    print(f"Load: {calculate_load(grid)}")


def cycle(grid: list[bytearray]):
    for _ in range(4):
        tilt_north(grid)
        rotate_grid(grid)


def find_cycle(grid: list[bytearray]) -> tuple[int, int]:
    cache = {grid_to_string(grid): [0]}
    index = 0
    while True:
        cycle(grid)
        indices = cache.get(grid_to_string(grid), [])
        indices.append(index)
        cache[grid_to_string(grid)] = indices
        if len(indices) == 2:
            return indices[1] - indices[0], (index + 1)
        index += 1


def solution_part_2():
    grid = load_data("input.txt")
    remaining_cycles = 1_000_000_000
    cycle_len, current_cycles = find_cycle(grid)
    remaining_cycles -= current_cycles
    remaining_cycles -= (remaining_cycles // cycle_len) * cycle_len
    for _ in range(remaining_cycles):
        cycle(grid)
    print(f"Load: {calculate_load(grid)}")


def main():
    solution_part_1()
    solution_part_2()


if __name__ == "__main__":
    main()
