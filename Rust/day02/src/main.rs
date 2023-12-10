use std::collections::HashMap;
use std::fs;
use lazy_static::lazy_static;

fn parse_game_data(data: &String) -> HashMap<i32, Vec<HashMap<String, i32>>> {
    let mut games : HashMap<i32, Vec<HashMap<String, i32>>> = HashMap::new();
    for line in data.lines() {
        // split line
        let split = line.split(":").collect::<Vec<&str>>();
        let game_index = split[0].split(" ").collect::<Vec<&str>>()[1].parse::<i32>().unwrap();
        let game_data_split = split[1].split(";").collect::<Vec<&str>>();
        let mut game = Vec::new();
        for single_draw in game_data_split {
            let mut draw = HashMap::new();
            let single_draw_split = single_draw.split(",").collect::<Vec<&str>>();
            for color in single_draw_split {
                let color_split = color.trim().split(" ").collect::<Vec<&str>>();
                let color_name = String::from(color_split[1]);
                let color_count = color_split[0].parse::<i32>().unwrap();
                draw.insert(color_name, color_count);
            }
            game.push(draw);
        }
        games.insert(game_index, game);
    }
    return games
}

lazy_static! {
    static ref GIVEN_CUBES: HashMap<String, i32> = {
        let mut m = HashMap::new();
        m.insert("red".to_string(), 12);
        m.insert("green".to_string(), 13);
        m.insert("blue".to_string(), 14);
        m
    };
}

fn test_game(game: &Vec<HashMap<String, i32>>) -> bool {
    for draw in game {
        for (color, count) in draw {
            if GIVEN_CUBES.contains_key(color) {
                let current_count = GIVEN_CUBES.get(color).unwrap();
                if current_count < count {
                    return false
                }
            } else {
                return false
            }
        }
    }
    return true
}
fn solution_part_1() {
    let data = fs::read_to_string("input.txt").expect("Unable to read file");
    let games = parse_game_data(&data);
    let mut sum = 0;
    for (game_index, game_data) in games {
        let result = test_game(&game_data);
        if result {
            sum += game_index
        }
    }
    println!("Sum: {}", sum);
}

fn get_power_of_game(game: &Vec<HashMap<String, i32>>) -> i32 {
    let mut current_cubes: HashMap<String, i32> = HashMap::new();
    current_cubes.insert("red".to_string(), 0);
    current_cubes.insert("green".to_string(), 0);
    current_cubes.insert("blue".to_string(), 0);
    for draw in game {
        for (color, count) in draw {
            let current_count = current_cubes.get_mut(color).unwrap();
            if (*current_count) < (*count) {
                *current_count = *count
            }
        }
    }
    return current_cubes.get("red").unwrap() * current_cubes.get("green").unwrap() * current_cubes.get("blue").unwrap()
}
fn solution_part_2() {
    let data = fs::read_to_string("input.txt").expect("Unable to read file");
    let games = parse_game_data(&data);
    let mut sum = 0;
    for (_, game_data) in games {
        sum += get_power_of_game(&game_data);
    }
    println!("Sum: {}", sum);
}
fn main() {
    solution_part_1();
    solution_part_2();
}
