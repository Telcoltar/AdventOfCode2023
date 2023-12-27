
#[derive(Debug, Clone, PartialEq, Eq)]
enum Condition {
    Operational,
    Damaged,
    Unknown
}

struct Process {
    conditions: Vec<Condition>,
    blocks: Vec<usize>,
    cache: Vec<i64>,
    cache_hits: i64,
    invocations: i64
}

const GRID_OVERSIZE: usize = 2;

impl Process {
    fn wrap_process_condition(&mut self, index: usize, block_index: usize) -> i64 {
        self.invocations += 1;
        let cache_value = self.cache[index * (self.blocks.len() + GRID_OVERSIZE) + block_index];
        if cache_value != -1 {
            self.cache_hits += 1;
        } else {
            let result = self.process_condition(index, block_index);
            self.cache[index * (self.blocks.len() + GRID_OVERSIZE) + block_index] = result;
        }
        self.cache[index * (self.blocks.len() + GRID_OVERSIZE) + block_index]
    }

    fn process_condition(&mut self, index: usize, block_index: usize) -> i64 {
        if index >= self.conditions.len() {
            return if block_index >= self.blocks.len() {
                1
            } else {
                0
            }
        }
        let current_condition = &self.conditions[index];
        match current_condition {
            Condition::Unknown => {
                let non_damaged_sum = self.wrap_process_condition(index + 1, block_index);
                let mut damaged_sum = 0;
                if block_index < self.blocks.len() && self.check_if_block_fits(self.blocks[block_index], index) {
                    damaged_sum = self.wrap_process_condition(index + self.blocks[block_index] + 1, block_index + 1);
                }
                return non_damaged_sum + damaged_sum;
            }
            Condition::Operational => {
                return self.wrap_process_condition(index + 1, block_index);
            }
            Condition::Damaged => {
                if block_index < self.blocks.len() && self.check_if_block_fits(self.blocks[block_index], index) {
                    return self.wrap_process_condition(index + self.blocks[block_index] + 1, block_index + 1);
                }
            }
        }
        0
    }

    fn check_if_block_fits(&self, block: usize, index: usize) -> bool {
        let mut current_index = index;
        let mut current_block = block;
        while current_index < self.conditions.len() && current_block > 0 &&
            (self.conditions[current_index] == Condition::Unknown || self.conditions[current_index] == Condition::Damaged) {
            current_index += 1;
            current_block -= 1;
        }
        current_block == 0 && (current_index == self.conditions.len() ||
            (self.conditions[current_index] == Condition::Unknown || self.conditions[current_index] == Condition::Operational))
    }

}

fn load_data(filename: &str) -> Vec<(Vec<Condition>, Vec<usize>)> {
    let data = std::fs::read_to_string(filename).expect("Error reading input file!");
    let mut result = Vec::new();
    for line in data.lines() {
        let mut conditions = Vec::new();
        let mut line_split = line.split(' ');
        for char in line_split.next().unwrap().chars() {
            match char {
                '.' => conditions.push(Condition::Operational),
                '#' => conditions.push(Condition::Damaged),
                '?' => conditions.push(Condition::Unknown),
                _ => panic!("Invalid condition!")
            }
        }
        let values: Vec<usize> = line_split.next().unwrap().trim().split(',').map(|x| x.parse().unwrap()).collect();
        result.push((conditions, values));
    }
    result
}
fn multiply_data(conditions: Vec<Condition>, values: Vec<usize>, factor: usize) -> (Vec<Condition>, Vec<usize>) {
    let mut result_conditions = conditions.to_vec();
    let mut result_values = values.to_vec();
    for _i in 0..(factor - 1) {
        result_conditions.push(Condition::Unknown);
        result_conditions.extend(conditions.to_vec());
        result_values.extend(values.to_vec());
    }
    (result_conditions, result_values)
}

fn solution(factor: usize) {
    let data = load_data("input.txt");

    let mut sum: i64 = 0;
    let mut cache_hits: i64 = 0;
    let mut invocations: i64 = 0;
    for (conditions, values) in data[5..].iter().cloned() {
        let (conditions, values) = multiply_data(conditions, values, factor);
        let cache = vec![-1; (conditions.len() + GRID_OVERSIZE) * (values.len() + GRID_OVERSIZE)];
        let mut process = Process { conditions, blocks: values, cache, cache_hits: 0, invocations: 0};
        sum += process.process_condition(0, 0);
        cache_hits += process.cache_hits;
        invocations += process.invocations;
    }
    println!("{}", sum);
    println!("Cache hits: {}", cache_hits);
    println!("Invocations: {}", invocations);

}

fn solution_part_1() {
    solution(1);
}

fn solution_part_2() {
    solution(5);
}

fn main() {
    solution_part_1();
    solution_part_2();
}
