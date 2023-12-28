use std::fs;

fn load_data(filename: &str) -> Vec<Vec<String>> {
    let mut data = Vec::new();
    let file = fs::read_to_string(filename).expect("Failed to read file");
    let mut current_map = Vec::new();
    for line in file.lines() {
        if line.is_empty() {
            data.push(current_map);
            current_map = Vec::new();
            continue;
        }
        current_map.push(line.to_string());
    }
    data.push(current_map);
    data
}

fn transpose_map(map: &Vec<String>) -> Vec<String> {
    let mut transposed_map = Vec::new();
    for i in 0..map[0].len() {
        let mut line = String::new();
        for j in 0..map.len() {
            line.push(map[j].chars().nth(i).unwrap());
        }
        transposed_map.push(line);
    }
    transposed_map
}

fn search_horizontal_reflection_line(map: &Vec<String>, equal_fn: fn(&str, &str) -> bool) -> Vec<usize> {
    let start_points = get_starting_points(map, equal_fn);
    let mut valid_points: Vec<usize> = Vec::new();
    for point in start_points {
        let mut lower_bound = point.0;
        let mut upper_bound = point.1;
        while lower_bound > 0 && upper_bound < map.len() - 1 && equal_fn(&map[lower_bound], &map[upper_bound]) {
            lower_bound -= 1;
            upper_bound += 1;
        }
        if (lower_bound == 0 || upper_bound == map.len() - 1) && equal_fn(&map[lower_bound], &map[upper_bound]) {
            valid_points.push(point.0);
        }
    }
    valid_points
}

fn search_horizontal_reflection_line_single(map: &Vec<String>) -> Option<usize> {
    let valid_points = search_horizontal_reflection_line(map, is_equal);
    if valid_points.len() == 1 {
        Some(valid_points[0])
    } else {
        None
    }
}
fn is_similar(line1: &str, line2: &str) -> bool {
    let mut diff = 0;
    for i in 0..line1.len() {
        if line1.chars().nth(i).unwrap() != line2.chars().nth(i).unwrap() {
            diff += 1;
        }
    }
    diff == 1 || diff == 0
}

fn is_equal(line1: &str, line2: &str) -> bool {
    line1 == line2
}

fn get_starting_points(map: &[String], equal_fn: fn(&str, &str) -> bool) -> Vec<(usize, usize)> {
    let mut start_points = Vec::new();
    for i in 0..map.len() - 1 {
        if equal_fn(&map[i], &map[i + 1]) {
            start_points.push((i, i + 1));
        }
    }
    start_points
}

fn solution_part_1() {
    let maps = load_data("input.txt");
    let mut vertical_sum = 0;
    let mut horizontal_sum = 0;
    for map in maps {
        if let Some(value) = search_horizontal_reflection_line_single(&map) {
            horizontal_sum += value + 1;
        } else if let Some(value) = search_horizontal_reflection_line_single(&transpose_map(&map)) {
            vertical_sum += value + 1;
        } else {
            println!("No solution found");
        }
    }
    let total_sum = vertical_sum + horizontal_sum * 100;
    println!("Total sum: {}", total_sum);
}

fn solution_part_2() {
    let maps = load_data("input.txt");
    let mut vertical_sum = 0;
    let mut horizontal_sum = 0;
    for map in maps {
        let strict_horizontal_points = search_horizontal_reflection_line(&map, is_equal);
        let horizontal_points = search_horizontal_reflection_line(&map, is_similar);
        let new_horizontal_points = horizontal_points.iter().filter(|x| !strict_horizontal_points.contains(x)).collect::<Vec<_>>();
        let strict_vertical_points = search_horizontal_reflection_line(&transpose_map(&map), is_equal);
        let vertical_points = search_horizontal_reflection_line(&transpose_map(&map), is_similar);
        let new_vertical_points = vertical_points.iter().filter(|x| !strict_vertical_points.contains(x)).collect::<Vec<_>>();
        if new_horizontal_points.len() == 1 {
            horizontal_sum += new_horizontal_points[0] + 1;
        } else if new_vertical_points.len() == 1 {
            vertical_sum += new_vertical_points[0] + 1;
        } else {
            println!("No solution found");
        }
    }
    let total_sum = vertical_sum + horizontal_sum * 100;
    println!("Total sum: {}", total_sum);
}
fn main() {
    solution_part_1();
    solution_part_2();
}
