import time
from enum import StrEnum, auto
from functools import lru_cache


class Condition(StrEnum):
    Operational = "."
    Damaged = "#"
    Unknown = "?"


def multiply_data(conditions: list[Condition], blocks: list[int], factor: int) -> tuple[list[Condition], list[int]]:
    new_conditions = conditions.copy()
    new_blocks = blocks.copy()
    for i in range(factor - 1):
        new_conditions.append(Condition.Unknown)
        new_conditions.extend(conditions)
        new_blocks.extend(blocks)
    return new_conditions, new_blocks


class Process:
    def __init__(self, conditions: list[Condition], blocks: list[int]):
        self.conditions = conditions
        self.blocks = blocks

    @lru_cache(maxsize=None)
    def process(self, cond_index: int, block_index: int) -> int:
        if cond_index >= len(self.conditions):
            if block_index >= len(self.blocks):
                return 1
            else:
                return 0
        condition = self.conditions[cond_index]
        if condition == Condition.Operational:
            return self.process(cond_index + 1, block_index)
        if condition == Condition.Damaged:
            if (block_index < len(self.blocks) and
                    self.check_if_block_fits(cond_index, self.blocks[block_index])):
                return self.process(cond_index + self.blocks[block_index] + 1, block_index + 1)
            return 0
        if condition == Condition.Unknown:
            undamaged_sum = self.process(cond_index + 1, block_index)
            damaged_sum = 0
            if block_index < len(self.blocks) and self.check_if_block_fits(cond_index, self.blocks[block_index]):
                damaged_sum = self.process(cond_index + self.blocks[block_index] + 1, block_index + 1)
            return undamaged_sum + damaged_sum

    def check_if_block_fits(self, cond_index: int, block: int) -> bool:
        current_cond_index = cond_index
        current_block = block
        while (current_cond_index < len(self.conditions) and current_block > 0 and
               (self.conditions[current_cond_index] == Condition.Damaged or
                self.conditions[current_cond_index] == Condition.Unknown)):
            current_cond_index += 1
            current_block -= 1
        return (current_block == 0 and
                (current_cond_index == len(self.conditions) or
                 self.conditions[current_cond_index] == Condition.Operational or
                 self.conditions[current_cond_index] == Condition.Unknown))


def conditions_to_str(conditions: list[Condition]) -> str:
    return "".join([c.value for c in conditions])


def load_data(filename: str) -> list[tuple[list[Condition], list[int]]]:
    data = []
    with open(filename, 'r') as file:
        for line in file:
            conditions, values = line.strip().split(' ')
            conditions = [Condition(condition.strip()) for condition in conditions]
            values = [int(value) for value in values.split(',')]
            data.append((conditions, values))
    return data


def solution(factor: int):
    data = load_data("input.txt")
    total_sum = 0
    for conditions, blocks in data:
        conditions, blocks = multiply_data(conditions, blocks, factor)
        process = Process(conditions, blocks)
        total_sum += process.process(0, 0)
    print(total_sum)


def solution_part_1():
    solution(1)


def solution_part_2():
    start = time.time_ns()
    solution(5)
    end = time.time_ns()
    print(f"Time: {(end - start)//1_000_000}ms")


def main():
    solution_part_1()
    solution_part_2()


if __name__ == "__main__":
    main()
