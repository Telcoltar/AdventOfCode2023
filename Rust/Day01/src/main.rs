use std::fs;
use phf::phf_map;
use regex::Regex;

static NUMBERS: phf::Map<&'static str, char> = phf_map! {
    "one" => '1',
    "two" => '2',
    "three" => '3',
    "four" => '4',
    "five" => '5',
    "six" => '6',
    "seven" => '7',
    "eight" => '8',
    "nine" => '9',
};

#[allow(dead_code)]
fn solution_part_1() {
    let content = fs::read_to_string("src/input.txt").expect("should work.");
    let mut line_numbers: Vec<i32> = Vec::new();
    for line in content.lines() {
        let mut line_number: [char; 2] = Default::default();
        for c in line.chars() {
            if c.is_digit(10) {
                line_number[0] = c;
                break;
            }
        }
        for c in line.chars().rev() {
            if c.is_digit(10) {
                line_number[1] = c;
                break;
            }
        }
        line_numbers.push(line_number.iter().collect::<String>().parse::<i32>().unwrap());
    }
    let sum: i32 = line_numbers.iter().sum();
    println!("{}", sum)
}

fn solution_part_2() {
    let content = fs::read_to_string("src/input.txt").expect("should work.");
    let mut line_numbers: Vec<i32> = Vec::new();
    let re = Regex::new(r"(one|two|three|four|five|six|seven|eight|nine|\d)").unwrap();
    let re_reverse_digits = Regex::new(r"(eno|owt|eerht|ruof|evif|xis|neves|thgie|enin|\d)").unwrap();
    for line in content.lines() {
        let mut line_number: [char; 2] = Default::default();
        let digit_match = re.find(line).unwrap();
        if digit_match.as_str().len() > 1 {
            line_number[0] = NUMBERS[digit_match.as_str()]
        } else {
            line_number[0] = digit_match.as_str().chars().next().unwrap();
        }
        let reverse_line: String = line.chars().rev().collect();
        let digit_reversed_match = re_reverse_digits.find(&reverse_line).unwrap();
        if digit_reversed_match.as_str().len() > 1 {
            let right_order_match: String =  digit_reversed_match.as_str().chars().rev().collect();
            line_number[1] = NUMBERS[&right_order_match];
        } else {
            line_number[1] = digit_reversed_match.as_str().chars().next().unwrap();
        }
        line_numbers.push(line_number.iter().collect::<String>().parse::<i32>().unwrap());
    }
    let sum: i32 = line_numbers.iter().sum();
    println!("{}", sum)
}

fn main() {
    solution_part_2()
}
