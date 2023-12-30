use std::collections::HashMap;
use std::fs;

#[derive(Debug, Clone)]
struct Box {
    index_map: HashMap<String, usize>,
    labels: Vec<i32>,
    current_index: usize,
}

fn load_data(filename: &str) -> Vec<String> {
    let file = fs::read_to_string(filename).expect("Failed to read file!");
    file.lines().next().unwrap().split(',').map(|s| s.to_string()).collect()
}

fn hash(input: &str) -> usize {
    let mut hash = 0;
    for c in input.chars() {
        hash += c as usize;
        hash *= 17;
        hash %= 256;
    }
    hash
}

fn solution_part_1() {
    let data = load_data("input.txt");
    let sum = data.iter().map(|s| hash(s)).sum::<usize>();
    println!("Sum: {}", sum);
}

fn solution_part_2() {
    let data = load_data("input.txt");
    let re_matcher = regex::Regex::new(r"([a-z]+)([=-])(\d*)").unwrap();
    let mut operations = Vec::new();
    for op in data {
        let caps = re_matcher.captures(&op).unwrap();
        let label = caps.get(1).unwrap().as_str().to_string();
        let op = caps.get(2).unwrap().as_str().chars().next().unwrap();
        operations.push((op, label.clone(), hash(&label), caps.get(3).unwrap().as_str().to_string()));
    }
    let mut boxes = vec![Box { index_map: HashMap::new(), labels: Vec::new(), current_index: 0 }; 256];
    for (op, label, hash, number) in operations {
        let current_box = &mut boxes[hash];
        if op == '=' {
            if let Some(index) = current_box.index_map.get(&label) {
                current_box.labels[*index] = number.parse::<i32>().unwrap();
            } else {
                current_box.index_map.insert(label.clone(), current_box.current_index);
                current_box.labels.push(number.parse::<i32>().unwrap());
                current_box.current_index += 1;
            }
        } else { // op == '-'
            if let Some(index) = current_box.index_map.get(&label) {
                current_box.labels[*index] = -1;
                current_box.index_map.remove(&label);
            }
        }
    }
    let mut total_sum = 0;
    for (index, b) in boxes.into_iter().enumerate() {
        if b.current_index > 0 {
            let mut box_sum = 0;
            let mut current_index = 1;
            for l in b.labels {
                if l != -1 {
                    box_sum += l * current_index;
                    current_index += 1;
                }
            }
            total_sum += box_sum * (index as i32 + 1);
        }
    }
    println!("Total sum: {}", total_sum);
}
fn main() {
    solution_part_1();
    solution_part_2();
}
