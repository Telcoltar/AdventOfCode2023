from functools import reduce


def load_data(filename: str) -> list[list[int]]:
    result = []
    with open(filename) as f:
        for line in f.readlines():
            result.append(list(map(int, line.split())))
    return result


def process_line(line: list[int]) -> list[list[int]]:
    result: list[list[int]] = [line]
    while True:
        current_line = []
        for i in range(len(result[-1]) - 1):
            current_line.append(result[-1][i + 1] - result[-1][i])
        result.append(current_line)
        if not any(current_line):
            break
    return result


def predict_end(lines: list[list[int]]) -> int:
    return sum(map(lambda line: line[-1], lines))


def predict_beginning(lines: list[list[int]]) -> int:
    return reduce(lambda acc, element: element - acc, map(lambda line: line[0], reversed(lines)))


def solution_part_1():
    data = load_data("input.txt")
    print(sum(map(predict_end, map(process_line, data))))


def solution_part_2():
    data = load_data("input.txt")
    print(sum(map(predict_beginning, map(process_line, data))))


def main():
    solution_part_1()
    solution_part_2()


if __name__ == "__main__":
    main()
