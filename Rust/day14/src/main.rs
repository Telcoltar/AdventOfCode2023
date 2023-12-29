use std::collections::{HashMap};

fn load_data(filename: &str) -> Vec<u8> {
    let data = std::fs::read_to_string(filename).unwrap();
    let mut platform = Vec::new();
    for line in data.lines() {
        let mut row = Vec::new();
        for c in line.chars() {
            match c {
                '.' => row.push(0),
                '#' => row.push(1),
                'O' => row.push(2),
                _ => panic!("Unknown field type")
            }
        }
        platform.push(row);
    }
    platform.iter().flat_map(|row| row.iter()).copied().collect()
}

fn roll_north(platform: &mut [u8], dimension: usize) {
    for y in 0..dimension {
        for x in 0..dimension {
            if platform[y * dimension + x] == 2 {
                let mut current_row = y;
                while current_row > 0 && platform[(current_row - 1) * dimension + x] == 0 {
                    current_row -= 1;
                }
                if current_row != y {
                    platform[current_row * dimension + x] = 2;
                    platform[y* dimension + x] = 0;
                }
            }
        }
    }
}

fn calculate_load(platform: &[u8], dimension: usize) -> usize {
    let mut load = 0;
    for y in 0..dimension {
        for x in 0..dimension {
            if platform[y * dimension + x] == 2 {
                load += dimension - y;
            }
        }
    }
    load
}

#[allow(dead_code)]
fn print_field(platform: &[u8], dimension: usize) {
    for y in 0..dimension {
        for x in 0..dimension {
            print!("{:?}", platform[y * dimension + x]);
        }
        println!();
    }
}

fn rotate_field(platform: &mut [u8], dimension: usize) {
    for i in 0..dimension {
        for j in i..dimension {
            platform.swap(i * dimension + j, j * dimension + i);
        }
    }
    for i in 0..dimension {
        for j in 0..dimension / 2 {
            platform.swap(i * dimension + j, i * dimension + dimension - j - 1);
        }
    }
}

fn solution_part_1() {
    let mut field = load_data("input.txt");
    let dimension = (field.len() as f64).sqrt() as usize;
    roll_north(&mut field, dimension);
    let load = calculate_load(&field, dimension);
    println!("Load: {}", load);
}

fn cycle_field(field: &mut Vec<u8>, dimension: usize) {
    for _ in 0..4 {
        roll_north(field, dimension);
        rotate_field(field, dimension);
    }
}

fn find_cycle(field: &mut Vec<u8>, dimension: usize) -> (usize, usize) {
    let mut cache: HashMap<Vec<u8>, Vec<usize>> = HashMap::new();
    cache.insert(field.clone(), vec![0]);
    for i in 0.. {
        cycle_field(field, dimension);
        let indices = cache.entry(field.clone()).or_default();
        indices.push(i + 1);
        if indices.len() == 2 {
            return (indices[1] - indices[0], i + 1);
        }
    }
    unreachable!();
}

fn solution_part_2() {
    let mut field = load_data("input.txt");
    let dimension = (field.len() as f64).sqrt() as usize;
    let mut remaining_cycles = 1000000000;
    let (cycle, current_cycles) = find_cycle(&mut field, dimension);
    remaining_cycles -= current_cycles;
    let cycles = remaining_cycles / cycle;
    remaining_cycles -= cycles * cycle;
    for _ in 0..remaining_cycles {
        cycle_field(&mut field, dimension);
    }
    let load = calculate_load(&field, dimension);
    println!("Load: {}", load);
}
fn main() {
    solution_part_1();
    solution_part_2();
}
