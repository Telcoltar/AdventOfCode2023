use std::collections::{HashSet, VecDeque};
use std::fs;
use std::ops::Index;

fn load_data(filename: &str) -> Vec<Vec<char>> {
    let mut data = Vec::new();
    let file = fs::read_to_string(filename).expect("Failed to read file!");
    for line in file.lines() {
        data.push(line.chars().collect::<Vec<char>>());
    }
    data
}

#[derive(Debug, Clone, Ord, PartialOrd, Eq, PartialEq, Hash)]
struct Point {
    x: usize,
    y: usize,
}

impl Point {
    fn move_point(&mut self, direction: u8) {
        match direction {
            0 => self.y -= 1,
            1 => self.x += 1,
            2 => self.y += 1,
            3 => self.x -= 1,
            _ => panic!("Invalid direction!"),
        }
    }
}

type Grid = Vec<Vec<char>>;

impl Index<&Point> for Grid {
    type Output = char;

    fn index(&self, index: &Point) -> &Self::Output {
        &self[index.y][index.x]
    }
}

fn pad_grid(grid: &mut Grid) {
    let mut new_grid = vec![vec!['#'; grid[0].len() + 2]; grid.len() + 2];
    for y in 0..grid.len() {
        new_grid[y + 1][0] = '#';
        for x in 0..grid[0].len() {
            new_grid[y + 1][x + 1] = grid[y][x];
        }
        new_grid[y + 1][grid[0].len() + 1] = '#';
    }
    *grid = new_grid;
}

#[derive(Debug, Clone, Ord, PartialOrd, Eq, PartialEq, Hash)]
struct Status {
    position: Point,
    direction: u8, // 0 = up, 1 = right, 2 = down, 3 = left
}

fn get_opposite_direction(direction: u8) -> u8 {
    match direction {
        0 => 2,
        1 => 3,
        2 => 0,
        3 => 1,
        _ => panic!("Invalid direction!"),
    }
}

impl Status {
    fn move_point(&mut self) {
        match self.direction {
            0 => self.position.y -= 1,
            1 => self.position.x += 1,
            2 => self.position.y += 1,
            3 => self.position.x -= 1,
            _ => panic!("Invalid direction!"),
        }
    }

}

fn solution_part_1() {
    let mut grid = load_data("input.txt");
    pad_grid(&mut grid);
    let status = Status { position: Point { x: 0, y: 1 }, direction: 1 };
    let energized_tiles = calculate_energized_tiles(&grid, status);
    println!("Part 1: {}", energized_tiles);
}

fn solution_part_2() {
    let mut grid = load_data("input.txt");
    pad_grid(&mut grid);
    let dimension = grid.len();
    let mut energized_tiles = Vec::new();
    for i in 0..dimension - 2 {
        energized_tiles.push(
            calculate_energized_tiles(&grid,
                                      Status { position: Point { x: 0, y: i + 1 }, direction: 1 }));
        energized_tiles.push(
            calculate_energized_tiles(&grid,
                                      Status { position: Point { x: dimension - 1, y: i + 1 }, direction: 3 }));
        energized_tiles.push(
            calculate_energized_tiles(&grid,
                                      Status { position: Point { x: i + 1, y: 0 }, direction: 2 }));
        energized_tiles.push(
            calculate_energized_tiles(&grid,
                                      Status { position: Point { x: i + 1, y: dimension - 1 }, direction: 0 }));
    }
    let max = energized_tiles.iter().max().unwrap();
    println!("Part 2: {}", max);
}

fn calculate_energized_tiles(grid: & Vec<Vec<char>>, status: Status) -> usize {
    let mut queue = VecDeque::new();
    let mut visited = HashSet::new();
    let mut count_points = HashSet::new();
    queue.push_back(status);
    while !queue.is_empty() {
        let mut status = queue.pop_front().unwrap();
        if visited.contains(&status) {
            continue;
        }
        visited.insert(status.clone());
        count_points.insert(status.position.clone());
        status.move_point();
        match grid[&status.position] {
            '.' => {
                let mut current_point = status.position.clone();
                while grid[&current_point] == '.' {
                    count_points.insert(current_point.clone());
                    current_point.move_point(status.direction);
                }
                current_point.move_point(get_opposite_direction(status.direction));
                status.position = current_point;
                queue.push_back(status)
            },
            '#' => {
                continue;
            },
            '|' => {
                if status.direction == 0 || status.direction == 2 {
                    queue.push_back(status)
                } else {
                    let mut second_status = status.clone();
                    status.direction = 0;
                    second_status.direction = 2;
                    queue.push_back(status);
                    queue.push_back(second_status);
                }
            },
            '-' => {
                if status.direction == 1 || status.direction == 3 {
                    queue.push_back(status)
                } else {
                    let mut second_status = status.clone();
                    status.direction = 1;
                    second_status.direction = 3;
                    queue.push_back(status);
                    queue.push_back(second_status);
                }
            },
            '\\' => {
                match status.direction {
                    0 => status.direction = 3,
                    1 => status.direction = 2,
                    2 => status.direction = 1,
                    3 => status.direction = 0,
                    _ => panic!("Invalid direction!"),
                }
                queue.push_back(status);
            },
            '/' => {
                match status.direction {
                    0 => status.direction = 1,
                    1 => status.direction = 0,
                    2 => status.direction = 3,
                    3 => status.direction = 2,
                    _ => panic!("Invalid direction!"),
                }
                queue.push_back(status);
            },
            _ => {
                panic!("Invalid character!")
            }
        }
    }
    count_points.len() - 1
}

fn main() {
    solution_part_1();
    let start = std::time::Instant::now();
    solution_part_2();
    println!("Time: {:?}", start.elapsed());
}
