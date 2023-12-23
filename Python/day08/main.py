import math
import re
from enum import StrEnum, auto
from typing import Callable


class Direction(StrEnum):
    R = "R"
    L = "L"


def load_data(filename: str) -> tuple[list[Direction], dict[str, dict[Direction, str]]]:
    with open(filename, encoding="utf-8") as f:
        directions = [Direction(d) for d in next(f).strip()]
        next(f)
        network = {}
        matcher = re.compile(r"\((\w+), (\w+)\)")
        for line in f:
            source, destinations = line.strip().split(" = ")
            destinations = matcher.fullmatch(destinations)
            network[source] = {Direction.L: destinations.group(1), Direction.R: destinations.group(2)}
        return directions, network


def follow_path(start: str, instruction: list[Direction], network: dict[str, dict[Direction, str]], check_fn: Callable[[str], bool]) -> int:
    current = start
    instruction_index = 0
    step = 0
    while check_fn(current):
        current = network[current][instruction[instruction_index]]
        step += 1
        instruction_index += 1
        if instruction_index == len(instruction):
            instruction_index = 0
    return step


def solution_part_1():
    directions, network = load_data("input.txt")
    steps = follow_path("AAA", directions, network, check_fn=lambda x: x != "ZZZ")
    print(steps)


def solution_part_2():
    directions, network = load_data("input.txt")
    cycles = []
    for node in network.keys():
        if node.endswith("A"):
            cycles.append(follow_path(node, directions, network, check_fn=lambda x: not x.endswith("Z")))
    kgv = math.lcm(*cycles)
    print(kgv)


def main():
    solution_part_1()
    solution_part_2()


if __name__ == "__main__":
    main()