use std::fs;

fn load_data_part_1(filename: &str) -> (Vec<f64>, Vec<f64>) {
    let data = fs::read_to_string(filename).expect("Unable to read file");
    let times = data.lines().nth(0).unwrap()
        .split(":").nth(1).unwrap().trim().split_whitespace()
        .map(|x| x.parse::<f64>().unwrap()).collect::<Vec<f64>>();
    let distances = data.lines().nth(1).unwrap()
        .split(":").nth(1).unwrap().trim().split_whitespace()
        .map(|x| x.parse::<f64>().unwrap()).collect::<Vec<f64>>();
    return (times, distances)
}

fn load_data_part_2(filename: &str) -> (f64, f64) {
    let data = fs::read_to_string(filename).expect("Unable to read file");
    let time = data.lines().nth(0).unwrap()
        .split(":").nth(1).unwrap().trim().replace(" ", "").parse::<f64>().unwrap();
    let distance = data.lines().nth(1).unwrap()
        .split(":").nth(1).unwrap().trim().replace(" ", "").parse::<f64>().unwrap();
    return (time, distance)
}

fn solution_part_1() {
    let (times, distances) = load_data_part_1("input.txt");
    let mut ranges: Vec<f64> = Vec::new();
    for (time, distance) in times.iter().zip(distances.iter()) {
        let range = calculate_range(time, distance);
        ranges.push(range);
    }
    let product = ranges.into_iter().reduce(|x, y| x * y).unwrap();
    println!("Product: {}", product);
}

fn solution_part_2() {
    let (time, distance) = load_data_part_2("input.txt");
    let range = calculate_range(&time, &distance);
    println!("Range: {}", range);
}

fn calculate_range(time: &f64, distance: &f64) -> f64 {
    let zero_one = time / 2.0 + ((time / 2.0).powf(2.0) - distance).sqrt();
    let zero_two = time / 2.0 - ((time / 2.0).powf(2.0) - distance).sqrt();
    let mut high = zero_one.floor();
    if high == zero_one {
        high -= 1.0;
    }
    let mut low = zero_two.ceil();
    if low == zero_two {
        low += 1.0;
    }
    let range = high - low + 1.0;
    range
}

fn main() {
    solution_part_1();
    solution_part_2();
}
