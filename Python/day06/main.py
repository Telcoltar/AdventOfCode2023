import math
from functools import reduce


def load_data_part_1(filename: str) -> tuple[list[int], list[int]]:
    with open(filename) as f:
        times = list(map(int, f.readline().strip().split(":")[1].strip().split()))
        distances = list(map(int, f.readline().strip().split(":")[1].strip().split()))
        return times, distances


def load_data_part_2(filename: str) -> tuple[int, int]:
    with open(filename) as f:
        time = int(f.readline().split(":")[1].strip().replace(" ", ""))
        distance = int(f.readline().split(":")[1].strip().replace(" ", ""))
        return time, distance


def solution_part_1():
    times, distances = load_data_part_1("input.txt")
    ranges: list[int] = []
    for time, distance in zip(times, distances):
        ranges.append(calculate_range(distance, time))
    print(reduce(lambda x,y: x*y, ranges, 1))


def solution_part_2():
    time, distance = load_data_part_2("input.txt")
    print(calculate_range(distance, time))


def calculate_range(distance, time):
    p_2 = time / 2
    d = math.sqrt(math.pow(time / 2.0, 2.0) - distance)
    zero_1 = p_2 - d
    zero_2 = p_2 + d
    high = math.floor(zero_2)
    low = math.ceil(zero_1)
    if high == zero_2:
        high -= 1
    if low == zero_1:
        low += 1
    return high - low + 1


def main():
    solution_part_1()
    solution_part_2()

if __name__ == "__main__":
    main()