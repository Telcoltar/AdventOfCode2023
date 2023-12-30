import re


def load_data(filename: str) -> list[str]:
    with open(filename, encoding="utf-8") as file:
        return file.readline().split(",")


def hash_alg(input_str: str) -> int:
    hash_sum = 0
    for char in input_str:
        hash_sum += ord(char)
        hash_sum *= 17
        hash_sum %= 256
    return hash_sum


def solution_part_1():
    ops = load_data("input.txt")
    total_sum = sum([hash_alg(op) for op in ops])
    print(total_sum)


class Box:

    def __init__(self):
        self.index_map: dict[str, int] = {}
        self.values: list[int] = []
        self.current_index = 0

    def insert(self, label: str, value: int) -> None:
        if label in self.index_map:
            self.values[self.index_map[label]] = value
        else:
            self.index_map[label] = self.current_index
            self.values.append(value)
            self.current_index += 1

    def remove(self, label):
        if label in self.index_map:
            self.values[self.index_map[label]] = -1
            del self.index_map[label]


def solution_part_2():
    ops = load_data("input.txt")
    boxes = [Box() for _ in range(256)]

    op_matcher = re.compile(r"([a-z]+)([=-])(\d*)")
    for op in ops:
        match = op_matcher.match(op)
        if match:
            label = match.group(1)
            operator = match.group(2)
            if operator == "=":
                boxes[hash_alg(label)].insert(label, int(match.group(3)))
            elif operator == "-":
                boxes[hash_alg(label)].remove(label)

    total_sum = 0
    for (index, b) in enumerate(boxes):
        if b.current_index > 0:
            box_sum = 0
            current_index = 1
            for value in b.values:
                if value != -1:
                    box_sum += (value * current_index)
                    current_index += 1
            total_sum += (box_sum * (index + 1))
    print(total_sum)


def main():
    solution_part_1()
    solution_part_2()


if __name__ == "__main__":
    main()