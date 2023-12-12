use std::fs;
use std::collections::HashSet;

struct Game {
    game_index: usize,
    winning_numbers: HashSet<usize>,
    my_numbers: HashSet<usize>,
}

fn parse_number_list(number_list: &str) -> HashSet<usize> {
    return number_list.split_whitespace().map(|x| x.parse::<usize>().unwrap()).collect();
}
fn load_data(filename: &str) -> Vec<Game> {
    let data = fs::read_to_string(filename).expect("Unable to read file");
    let mut games: Vec<Game> = Vec::new();
    for line in data.lines() {
        let mut line_split = line.split(":");
        let game_split = line_split.next().unwrap().split_whitespace().collect::<Vec<&str>>();
        let game_index = game_split[1].parse::<usize>().unwrap();
        let mut number_split = line_split.next().unwrap().split("|");
        let winning_numbers = parse_number_list(number_split.next().unwrap());
        let my_numbers = parse_number_list(number_split.next().unwrap());
        games.push(Game { game_index, winning_numbers, my_numbers });
    }
    return games
}

fn solution_part_1() {
    let games = load_data("input.txt");
    let mut total_score = 0;
    for game in games {
        let mut matches = 0;
        for number in game.my_numbers {
            if game.winning_numbers.contains(&number) {
                matches += 1;
            }
        }
        if matches > 0 {
            total_score += 2_i32.pow(matches as u32 - 1);
        }
    }
    println!("Total score: {}", total_score);
}

fn solution_part_2() {
    let games = load_data("input.txt");
    let mut card_pile = vec![1; games.len()];
    for game in games {
        let mut matches = 0;
        for number in game.my_numbers {
            if game.winning_numbers.contains(&number) {
                matches += 1;
            }
        }
        for index in game.game_index..game.game_index + matches {
            card_pile[index] += card_pile[game.game_index - 1];
        }
    }
    println!("Total score: {}", card_pile.into_iter().sum::<usize>());
}
fn main() {
    solution_part_1();
    solution_part_2();
}
