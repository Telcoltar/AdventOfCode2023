
fn load_data(filename: &str) -> Vec<Vec<i64>> {
    let data = std::fs::read_to_string(filename).unwrap();
    data.lines().map(|line| line.split_whitespace().map(|x| x.parse().unwrap()).collect()).collect()
}

fn process_line(line: Vec<i64>) -> Vec<Vec<i64>> {
    let mut result: Vec<Vec<i64>> = vec![line];
    // construct new line from differences
    loop {
        let mut new_line: Vec<i64> = Vec::new();
        for i in 0..result[result.len() - 1].len() - 1 {
            new_line.push(result[result.len() - 1][i + 1] - result[result.len() - 1][i]);
        }
        result.push(new_line.clone());
        if new_line.iter().all(|&x| x == 0) {
            break;
        }
    }
    result
}

fn predict_value_end(lines: Vec<Vec<i64>>) -> i64 {
    lines.iter().rev().fold(0, |acc, line| {
        acc + line[line.len() - 1]
    })
}

fn predict_value_beginning(lines: Vec<Vec<i64>>) -> i64 {
    lines.iter().rev().fold(0, |acc, line| {
        line[0] - acc
    })
}
fn solution_part_1() {
    let data = load_data("input.txt");
    let prediction = data.into_iter().map(process_line).map(predict_value_end).sum::<i64>();
    println!("{:?}", prediction)
}

fn solution_part_2() {
    let data = load_data("input.txt");
    let prediction = data.into_iter().map(process_line).map(predict_value_beginning).sum::<i64>();
    println!("{:?}", prediction)
}
fn main() {
    solution_part_1();
    solution_part_2();
}
