import re

DIGIT_WORDS = {
    "one": "1",
    "two": "2",
    "three": "3",
    "four": "4",
    "five": "5",
    "six": "6",
    "seven": "7",
    "eight": "8",
    "nine": "9"
}

def solution_part_1():
    with open("input.txt", "r") as f:
        line_numbers = []
        for line in f.readlines():
            line_number = ""
            for c in line:
                if c.isdigit():
                    line_number += c
                    break
            for c in line[::-1]:
                if c.isdigit():
                    line_number += c
                    break
            line_numbers.append(line_number)
        print(sum([int(line_number) for line_number in line_numbers]))


def solution_part_2():
    with open("input.txt", "r") as f:
        line_numbers = []
        for i, line in enumerate(f.readlines()):
            line_number = ""
            digits = re.compile(r"(one|two|three|four|five|six|seven|eight|nine|\d)")
            match = digits.search(line)
            if len(match[0]) == 1:
                line_number += match[0]
            else:
                line_number += DIGIT_WORDS[match[0]]
            digits_reverse = re.compile(r"(eno|owt|eerht|ruof|evif|xis|neves|thgie|enin|\d)")
            reverse_matche = digits_reverse.findall(line[::-1])
            if len(reverse_matche[0]) == 1:
                line_number += reverse_matche[0]
            else:
                line_number += DIGIT_WORDS[reverse_matche[0][::-1]]
            line_numbers.append(line_number)
        print(sum([int(line_number) for line_number in line_numbers]))

def main():
    solution_part_2()


if __name__ == "__main__":
    main()