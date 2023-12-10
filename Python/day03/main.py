from collections import defaultdict
from itertools import product


def load_board(filename: str) -> list[(int, list[((int, int), str)])]:
    numbers: list[(int, list[((int, int), str)])] = []
    padded_board: list[str] = []
    with open(filename) as f:
        first_line = f.readline().strip()
        board_width = len(first_line) + 2
        padded_board.append("." * board_width)
        padded_board.append("." + first_line + ".")
        for line in f:
            padded_board.append("." + line.strip() + ".")
        padded_board.append("." * board_width)
    for y, line in enumerate(padded_board):
        number_in_progress = False
        current_number = ""
        start: (int, int) = (0,0)
        for x, char in enumerate(line):
            if char.isdigit():
                if number_in_progress:
                    current_number += char
                else:
                    start = (x - 1, y - 1)
                    number_in_progress = True
                    current_number = char
            else:
                if number_in_progress:
                    number_in_progress = False
                    symbols: list[((int, int), str)] = []
                    for i in range(start[0], x + 1):
                        if padded_board[y - 1][i] != "." and not padded_board[y - 1][i].isdigit():
                            symbols.append(((i, y - 1), padded_board[y - 1][i]))
                    for i in range(start[0], x + 1):
                        if padded_board[y + 1][i] != "." and not padded_board[y + 1][i].isdigit():
                            symbols.append(((i, y + 1), padded_board[y + 1][i]))
                    if padded_board[y][start[0]] != "." and not padded_board[y][start[0]].isdigit():
                        symbols.append(((start[0], y), padded_board[y][start[0]]))
                    if padded_board[y][x] != "." and not padded_board[y][x].isdigit():
                        symbols.append(((x, y), padded_board[y][x]))
                    numbers.append((int(current_number), symbols))
    return numbers


def solution_part_1():
    board = load_board("input.txt")
    part_number_sum = 0
    for number in board:
        value, symbols = number
        if len(symbols) > 0:
            part_number_sum += value
    print(part_number_sum)


def solution_part_2():
    board = load_board("input.txt")
    numbers_for_stars: dict[(int, int), list[int]] = defaultdict(lambda: [])
    gear_ratio_sum = 0
    for number in board:
        value, symbols = number
        for symbol in symbols:
            if symbol[1] == "*":
                numbers_for_stars[symbol[0]].append(value)
    for star in numbers_for_stars:
        if len(numbers_for_stars[star]) == 2:
            gear_ratio_sum += numbers_for_stars[star][0] * numbers_for_stars[star][1]
    print(gear_ratio_sum)


def main():
    solution_part_1()
    solution_part_2()


if __name__ == "__main__":
    main()
