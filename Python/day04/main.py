def load_data(filename: str) -> list[tuple[set[int], set[int]]]:
    with open(filename) as f:
        cards: list[tuple[set[int], set[int]]] = []
        for line in f.readlines():
            _, numbers = line.split(":")
            winning_numbers, my_numbers = numbers.split("|")
            winning_numbers = set(map(int, winning_numbers.split()))
            my_numbers = set(map(int, my_numbers.split()))
            cards.append((winning_numbers, my_numbers))
    return cards


def solution_part_1():
    cards = load_data("input.txt")
    score = 0
    for card in cards:
        winning_numbers, my_numbers = card
        matches = len(winning_numbers.intersection(my_numbers))
        if matches > 0:
            score += pow(2, len(winning_numbers.intersection(my_numbers)) - 1)
    print(score)


def solution_part_2():
    cards = load_data("input.txt")
    pile : list[int] = [1 for _ in range(len(cards))]
    for index, card in enumerate(cards):
        winning_numbers, my_numbers = card
        matches = len(winning_numbers.intersection(my_numbers))
        for i in range(index + 1, index + matches + 1):
            pile[i] += pile[index]
    print(sum(pile))


def main():
    solution_part_1()
    solution_part_2()


if __name__ == "__main__":
    main()
