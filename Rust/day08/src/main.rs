use std::collections::HashMap;
use regex::Regex;
use num_integer::lcm;

#[derive(Debug)]
enum Direction {
    R,
    L,
}

fn load_data(filename: &str) -> (Vec<Direction>, HashMap<String, (String, String)>) {
    let binding = std::fs::read_to_string(filename).unwrap();
    let mut data = binding.lines();
    let mut instructions = Vec::new();
    for direction in data.next().unwrap().chars() {
        if direction == 'R' {
            instructions.push(Direction::R);
        } else {
            instructions.push(Direction::L);
        }
    }
    let mut network = HashMap::new();
    data.next();
    let re = Regex::new(r"\((\w{3}), (\w{3})\)").unwrap();
    for line in data {
        let mut line = line.split(" = ");
        let address = line.next().unwrap().to_string();
        let value = line.next().unwrap().to_string();
        match re.captures(&value) {
            Some(caps) => {
                let value = (caps[1].to_string(), caps[2].to_string());
                network.insert(address, value);
            }
            None => {
                let value = (value, "".to_string());
                network.insert(address, value);
            }
        }
    }
    return (instructions, network);
}

fn follow_path(starts: &str,
               instructions: &Vec<Direction>,
               network: &HashMap<String, (String, String)>,
               terminate_fn: &dyn Fn(&str) -> bool) -> i32 {
    let mut current = starts;
    let mut current_instruction_index = 0;
    let mut steps = 0;
    while terminate_fn(current) {
        steps += 1;
        match instructions[current_instruction_index] {
            Direction::L => {
                let (next, _) = network.get(current).unwrap();
                // println!("{} -> {}", current, next);
                current = next;
                current_instruction_index += 1;
            }
            Direction::R => {
                let (_, next) = network.get(current).unwrap();
                // println!("{} -> {}", current, next);
                current = next;
                current_instruction_index += 1;
            }
        }
        if current_instruction_index == instructions.len() {
            current_instruction_index = 0;
        }
    }
    return steps;
}
fn solution_part_1() {
    let (instruction, network) = load_data("input.txt");
    let steps = follow_path("AAA", &instruction, &network, &|current| current != "ZZZ");
    println!("Steps: {}", steps);
}

fn solution_part_2() {
    let (instruction, network) = load_data("input.txt");
    // find nodes ending in "A"
    let mut cycles: Vec<i64> = Vec::new();
    for node in network.keys() {
        if node.ends_with("A") {
            cycles.push(follow_path(node, &instruction, &network, &|current| !current.ends_with("Z")) as i64);
        }
    }
    let cycles = cycles.iter().fold(1, |acc, x| lcm(acc, *x));
    println!("Cycles: {:?}", cycles);
}
fn main() {
    solution_part_1();
    solution_part_2()
}
