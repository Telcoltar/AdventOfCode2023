from collections import defaultdict
from enum import IntEnum
from functools import cmp_to_key, partial, reduce
from typing import Callable


class TypeOfHand(IntEnum):
    FiveOfAKind = 7
    FourOfAKind = 6
    FullHouse = 5
    ThreeOfAKind = 4
    TwoPair = 3
    OnePair = 2
    HighCard = 1


def load_data(filename: str) -> list[tuple[str, int]]:
    result: list[tuple[str, int]] = []
    with open(filename) as f:
        for line in f.readlines():
            hand, bid = line.split(" ")
            result.append((hand, int(bid)))
    return result


def count_cards(hand: str) -> dict[str, int]:
    count: dict[str, int] = defaultdict(lambda: 0)
    for card in hand:
        count[card] += 1
    return count


def get_type_for_values(values: list[int]) -> TypeOfHand:
    if values[0] == 5:
        return TypeOfHand.FiveOfAKind
    if values[0] == 4:
        return TypeOfHand.FourOfAKind
    elif values == [3, 2]:
        return TypeOfHand.FullHouse
    elif values == [3, 1, 1]:
        return TypeOfHand.ThreeOfAKind
    elif values == [2, 2, 1]:
        return TypeOfHand.TwoPair
    elif values[0] == 2:
        return TypeOfHand.OnePair
    return TypeOfHand.HighCard


def parse_type_part_1(hand: str) -> TypeOfHand:
    count = count_cards(hand)
    values = list(count.values())
    values.sort(reverse=True)
    return get_type_for_values(values)


def parse_type_part_2(hand: str) -> TypeOfHand:
    count = count_cards(hand)
    joker_count = count["J"]
    if joker_count == 5 or joker_count == 4:
        return TypeOfHand.FiveOfAKind
    del count["J"]
    values = list(count.values())
    values.sort(reverse=True)
    values[0] += joker_count
    return get_type_for_values(values)


def card_value(card: str, j_value: int) -> int:
    if card == "A":
        return 14
    elif card == "K":
        return 13
    elif card == "Q":
        return 12
    elif card == "J":
        return j_value
    elif card == "T":
        return 10
    else:
        return int(card)


def compare_hands(hand_1: (str, TypeOfHand, int), hand_2: (str, TypeOfHand, int), j_value: int) -> int:
    if hand_1[1].value == hand_2[1].value:
        for card_1, card_2 in zip(hand_1[0], hand_2[0]):
            if card_value(card_1, j_value) > card_value(card_2, j_value):
                return 1
            elif card_value(card_1, j_value) < card_value(card_2, j_value):
                return -1
        return 0
    elif hand_2[1].value > hand_1[1].value:
        return -1
    return 1


def solution_part_1():
    print(solution("input.txt", 11, parse_type_part_1))


def solution_part_2():
    print(solution("input.txt", 1, parse_type_part_2))


def solution(filename: str, j_value: int, parse_fn: Callable[[str], TypeOfHand]) -> int:
    data = list(map(lambda x: (x[0], parse_fn(x[0]), x[1]), load_data(filename)))
    data.sort(key=cmp_to_key(partial(compare_hands, j_value=j_value)))
    return reduce(lambda acc, hand: acc + (hand[0] + 1) * hand[1][2], enumerate(data), 0)


def main():
    solution_part_1()
    solution_part_2()


if __name__ == "__main__":
    main()
