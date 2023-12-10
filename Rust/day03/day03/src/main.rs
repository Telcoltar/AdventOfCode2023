use std::collections::HashMap;
use std::fs;

fn load_board(filename: &str) -> (Vec<((usize, usize), (usize, usize), i32)>, Vec<(usize, usize)>) {
    let data = fs::read_to_string(filename).expect("Unable to read file");
    let mut board: Vec<((usize, usize), (usize, usize), i32)> = Vec::new();
    let mut symbols: Vec<(usize, usize)> = Vec::new();
    for (y, line) in data.lines().enumerate() {
        let mut number_parsing_in_progress = false;
        let mut current_number_str: Vec<char> = Vec::new();
        let mut starting_point: (usize, usize) = (0,0);
        for (x, c) in line.chars().enumerate() {
            if c.is_digit(10) {
                if number_parsing_in_progress {
                    current_number_str.push(c);
                } else {
                    number_parsing_in_progress = true;
                    starting_point = (x, y);
                    current_number_str = vec![c];
                }
            } else {
                if number_parsing_in_progress {
                    number_parsing_in_progress = false;
                    let current_number = current_number_str.iter().collect::<String>().parse::<i32>().unwrap();
                    board.push((starting_point, (x + 1,y + 2), current_number));
                }
                if c != '.' {
                    symbols.push((x + 1, y + 1));
                }
            }
        }
        // check if number parsing is still in progress
        if number_parsing_in_progress {
            let current_number = current_number_str.iter().collect::<String>().parse::<i32>().unwrap();
            board.push((starting_point, (line.len() + 1,y + 2), current_number));
        }
    }
    return (board, symbols)
}

fn load_board_clean(filename: &str) -> Vec<Vec<char>> {
    let data = fs::read_to_string(filename).expect("Unable to read file");
    let mut board: Vec<Vec<char>> = Vec::new();
    let first_line = data.lines().next().unwrap();
    // fill first line with empty chars
    let mut first_line_vec: Vec<char> = Vec::new();
    for _ in 0..first_line.len() + 2 {
        first_line_vec.push('.');
    }
    board.push(first_line_vec);
    for line in data.lines() {
        let mut line_vec: Vec<char> = Vec::new();
        line_vec.push('.');
        for c in line.chars() {
            line_vec.push(c);
        }
        line_vec.push('.');
        board.push(line_vec);
    }
    // fill last line with empty chars
    let mut last_line_vec: Vec<char> = Vec::new();
    for _ in 0..first_line.len() + 2 {
        last_line_vec.push('.');
    }
    board.push(last_line_vec);
    return board
}

fn print_sub_board(board: &Vec<Vec<char>>, start: (usize, usize), end: (usize, usize)) {
    for y in start.1..end.1 + 1 {
        for x in start.0..end.0 + 1 {
            print!("{}", board[y][x]);
        }
        println!();
    }
}

fn get_star(board: &Vec<Vec<char>>, start: (usize, usize), end: (usize, usize)) -> Option<(usize, usize)> {
    for y in start.1..end.1 + 1 {
        if board[y][start.0] == '*' {
            return Some((start.0, y))
        }
        if board[y][end.0] == '*' {
            return Some((end.0, y))
        }
    }
    for x in start.0..end.0 {
        if board[start.1][x] == '*'{
            return Some((x, start.1))
        }
        if board[end.1][x] == '*'{
            return Some((x, end.1))
        }
    }
    return None
}

fn solution_part_1() {
    let (numbers, symbols) = load_board("input.txt");
    let mut symbol_numbers_sum = 0;
    for (start, end, number) in numbers {
        for symbol in symbols.iter() {
            if symbol.0 >= start.0 && symbol.0 <= end.0 && symbol.1 >= start.1 && symbol.1 <= end.1 {
                symbol_numbers_sum += number;
                break
            }
        }
    }
    println!("Solution part 1: {}", symbol_numbers_sum);
}

fn solution_part_2(filename: &str) {
    let (numbers, _) = load_board(filename);
    let board = load_board_clean(filename);
    let mut star_numbers: HashMap<(usize, usize), Vec<i32>> = HashMap::new();
    for (start, end, number) in numbers {
        let star = get_star(&board, (start.0, start.1), (end.0, end.1));
        if star.is_some() {
            let star = star.unwrap();
            if star_numbers.contains_key(&star) {
                star_numbers.get_mut(&star).unwrap().push(number);
            } else {
                star_numbers.insert(star, vec![number]);
            }
        }
    }
    let mut gear_ratio_sum = 0;
    for star in star_numbers.keys() {
        if star_numbers.get(star).unwrap().len() == 2 {
            gear_ratio_sum += star_numbers.get(star).unwrap()[0] * star_numbers.get(star).unwrap()[1];
        }
    }
    println!("Solution part 2: {}", gear_ratio_sum);
}
fn main() {
    solution_part_1();
    solution_part_2("input.txt")
}
