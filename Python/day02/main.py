from enum import StrEnum, auto


class CubeColor(StrEnum):
    RED = auto()
    GREEN = auto()
    BLUE = auto()


MAX_CUBES = {
    CubeColor.RED: 12,
    CubeColor.GREEN: 13,
    CubeColor.BLUE: 14
}


def load_game(filename: str) -> dict[int, dict[CubeColor, int]]:
    game_data = {}
    with open(filename) as f:
        for line in f:
            game, draws = line.split(":")
            _, game_id = game.split(" ")
            game_data[int(game_id)] = {}
            cubes_max = {
                CubeColor.RED: 0,
                CubeColor.GREEN: 0,
                CubeColor.BLUE: 0
            }
            for draw in draws.split(";"):
                for cube in draw.split(","):
                    cube_count, cube_color = cube.strip().split(" ")
                    cubes_max[CubeColor(cube_color)] = max(cubes_max[CubeColor(cube_color)], int(cube_count))
            game_data[int(game_id)] = cubes_max
    return game_data


def solution_part_1(game_data: dict[int, dict[CubeColor, int]]):
    possible_games_sum = 0
    for index, game in game_data.items():
        if (
            game[CubeColor.RED] <= MAX_CUBES[CubeColor.RED] and
            game[CubeColor.GREEN] <= MAX_CUBES[CubeColor.GREEN] and
            game[CubeColor.BLUE] <= MAX_CUBES[CubeColor.BLUE]
        ):
            possible_games_sum += index
    print(possible_games_sum)


def solution_part_2(game_data: dict[int, dict[CubeColor, int]]):
    power_sum = 0
    for game in game_data.values():
        power_sum += game[CubeColor.RED] * game[CubeColor.GREEN] * game[CubeColor.BLUE]
    print(power_sum)


def main():
    solution_part_2(load_game("input.txt"))


if __name__ == "__main__":
    main()
