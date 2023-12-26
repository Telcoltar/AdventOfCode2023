Point = dict[str, int]


def load_data(filename: str) -> list[Point]:
    galaxies: list[Point] = []
    with open(filename, 'r') as file:
        for y, line in enumerate(file):
            for x, char in enumerate(line):
                if char == "#":
                    galaxies.append({"x": x, "y": y})
    return galaxies


def spread_galaxies(galaxies: list[Point], spread_factor: int):
    for coord in ["x", "y"]:
        galaxies.sort(key=lambda g: g[coord])
        current_spread = 0
        curren_pos = 0
        for galaxy in galaxies:
            diff = galaxy[coord] - curren_pos
            curren_pos = galaxy[coord]
            if diff > 1:
                current_spread += (diff - 1) * (spread_factor - 1)
            galaxy[coord] += current_spread


def solution(galaxies: list[Point], spread_factor: int):
    spread_galaxies(galaxies, spread_factor)
    total_distance = sum_distance(galaxies)
    print(total_distance)


def sum_distance(galaxies: list[Point]) -> int:
    total_distance = 0
    for i in range(len(galaxies)):
        for j in range(i + 1, len(galaxies)):
            total_distance += abs(galaxies[i]["x"] - galaxies[j]["x"]) + abs(galaxies[i]["y"] - galaxies[j]["y"])
    return total_distance


def solution_part_1():
    galaxies = load_data("input.txt")
    solution(galaxies, 2)


def solution_part_2():
    galaxies = load_data("input.txt")
    solution(galaxies, 1_000_000)


def main():
    solution_part_1()
    solution_part_2()


if __name__ == "__main__":
    main()
