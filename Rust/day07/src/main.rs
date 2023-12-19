use std::cmp::Ordering;
use std::collections::HashMap;

#[derive(Debug)]
enum TypeOfHand {
    FiveOfAKind,
    FourOfAKind,
    FullHouse,
    ThreeOfAKind,
    TwoPair,
    OnePair,
    HighCard,
}

impl TypeOfHand {
    fn value(&self) -> i32 {
        match self {
            TypeOfHand::FiveOfAKind => 7,
            TypeOfHand::FourOfAKind => 6,
            TypeOfHand::FullHouse => 5,
            TypeOfHand::ThreeOfAKind => 4,
            TypeOfHand::TwoPair => 3,
            TypeOfHand::OnePair => 2,
            TypeOfHand::HighCard => 1,
        }
    }
}

fn load_data(filename: &str) -> Vec<(String, i32)> {
    let data = std::fs::read_to_string(filename).unwrap();
    let mut result = Vec::new();
    for line in data.lines() {
        let mut split = line.split(" ");
        let hand = split.nth(0).unwrap().to_string();
        let bid: i32 = split.nth(0).unwrap().parse().unwrap();
        result.push((hand, bid));
    }
    return result;
}

fn count_cards(hand: &str) -> HashMap<char, i32> {
    let hand = hand.chars();
    let mut count = HashMap::new();
    for card in hand {
        let count = count.entry(card).or_insert(0);
        *count += 1;
    }
    return count;
}

fn parse_values(values: Vec<i32>) -> TypeOfHand {
    if values.contains(&5) {
        return TypeOfHand::FiveOfAKind;
    }
    if values.contains(&4) {
        return TypeOfHand::FourOfAKind;
    }
    if values.contains(&3) && values.contains(&2) {
        return TypeOfHand::FullHouse;
    }
    if values.contains(&3) {
        return TypeOfHand::ThreeOfAKind;
    }
    if values.contains(&2) && values.len() == 3 {
        return TypeOfHand::TwoPair;
    }
    if values.contains(&2) {
        return TypeOfHand::OnePair;
    }
    return TypeOfHand::HighCard;
}

fn parse_type_part_1(hand: &str) -> TypeOfHand {
    let count = count_cards(hand);
    let values = count.into_values().collect::<Vec<i32>>();
    return parse_values(values);
}

fn parse_type_part_2(hand: &str) -> TypeOfHand {
    let mut count = count_cards(hand);
    let joker_count = count.get(&'J').cloned().unwrap_or(0);
    if joker_count == 5 || joker_count == 4 {
        return TypeOfHand::FiveOfAKind;
    }
    count.remove(&'J');
    let mut values = count.into_values().collect::<Vec<i32>>();
    values.sort_unstable();
    values.reverse();
    values[0] += joker_count;
    return parse_values(values);
}

fn compare_hands(
    hand_1: (&String, &TypeOfHand),
    hand_2: (&String, &TypeOfHand),
    j_value: i32,
) -> Ordering {
    if hand_1.1.value() > hand_2.1.value() {
        return Ordering::Greater;
    } else if hand_1.1.value() < hand_2.1.value() {
        return Ordering::Less;
    } else {
        for (char_1, char_2) in hand_1.0.chars().zip(hand_2.0.chars()) {
            if char_value(char_1, j_value) > char_value(char_2, j_value) {
                return Ordering::Greater;
            } else if char_value(char_1, j_value) < char_value(char_2, j_value) {
                return Ordering::Less;
            }
        }
    }
    return Ordering::Equal;
}

fn char_value(char: char, j_value: i32) -> i32 {
    match char {
        'A' => 14,
        'K' => 13,
        'Q' => 12,
        'J' => j_value,
        'T' => 10,
        _ => char.to_digit(10).unwrap() as i32,
    }
}

fn solution(filename: &str, j_value: i32, parse_fn: fn(&str) -> TypeOfHand) {
    let data = load_data(filename);
    let mut data = data
        .into_iter()
        .map(|(hand, bid)| ((parse_fn(&hand), hand), bid))
        .collect::<Vec<((TypeOfHand, String), i32)>>();
    data.sort_unstable_by(|a, b| compare_hands((&a.0 .1, &a.0 .0), (&b.0 .1, &b.0 .0), j_value));
    let winnings = data.iter().enumerate().fold(0, |acc, (index, (_, bid))| {
        let index = index + 1;
        let bid = bid * index as i32;
        acc + bid
    });
    println!("Winnings: {}", winnings);
}
fn solution_part_1() {
    solution("input.txt", 11, parse_type_part_1)
}

fn solution_part_2() {
    solution("input.txt", 1, parse_type_part_2)
}

fn main() {
    solution_part_1();
    solution_part_2();
}
