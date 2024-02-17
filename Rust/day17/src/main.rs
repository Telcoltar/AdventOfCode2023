use std::cmp::Ordering;
use std::collections::{BinaryHeap, HashMap};

#[derive(Copy, Clone, Eq, PartialEq, Debug)]
enum Direction {
    Up = 0,
    Down = 1,
    Left = 2,
    Right = 3,
}

#[derive(Debug)]
struct Origin {
    direction: Direction,
    distance: usize,
}



fn load_data(filename: &str) -> Grid {
    let data = std::fs::read_to_string(filename).unwrap();
    return data.lines().map(|line| {
        line.chars().map(|c| c.to_digit(10).unwrap() as u8).collect::<Vec<u8>>()
    }).collect::<Vec<Vec<u8>>>();
}

type Grid = Vec<Vec<u8>>;

fn get_neighbours(grid: &Grid, point: &Point) -> Vec<Point> {
    let mut neighbours = Vec::new();
    if point.x > 0 { neighbours.push(Point { x: point.x - 1, y: point.y }); }
    if point.y > 0 { neighbours.push(Point { x: point.x, y: point.y - 1 }); }
    if point.x < grid.len() - 1 { neighbours.push(Point { x: point.x + 1, y: point.y }); }
    if point.y < grid.len() - 1 { neighbours.push(Point { x: point.x, y: point.y + 1 }); }
    neighbours
}

#[derive(Copy, Clone, Eq, PartialEq, Debug)]
struct Point {
    x: usize,
    y: usize,
}

impl Point {
    fn get_point_in_direction(&self, direction: Direction, grid_dim: usize) -> Option<Point> {
        match direction {
            Direction::Up => {
                if self.y > 0 {
                    Some(Point { x: self.x, y: self.y - 1 })
                } else {
                    None
                }
            },
            Direction::Down => {
                if self.y < grid_dim - 1 {
                    Some(Point { x: self.x, y: self.y + 1 })
                } else {
                    None
                }
            },
            Direction::Left => {
                if self.x > 0 {
                    Some(Point { x: self.x - 1, y: self.y })
                } else {
                    None
                }
            },
            Direction::Right => {
                if self.x < grid_dim - 1 {
                    Some(Point { x: self.x + 1, y: self.y })
                } else {
                    None
                }
            },
        }
    }

}

#[derive(Debug)]
struct State {
    point: Point,
    origin: Origin,
    cost: usize,
}

impl Eq for State {}

impl PartialEq<Self> for State {
    fn eq(&self, other: &Self) -> bool {
        self.cost == other.cost
    }
}

impl PartialOrd<Self> for State {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(other.cost.cmp(&self.cost))
    }
}

impl Ord for State {
    fn cmp(&self, other: &Self) -> Ordering {
        other.cost.cmp(&self.cost)
    }
}

fn get_orthogonal_directions(direction: Direction) -> Vec<Direction> {
    match direction {
        Direction::Up => vec![Direction::Left, Direction::Right],
        Direction::Down => vec![Direction::Left, Direction::Right],
        Direction::Left => vec![Direction::Up, Direction::Down],
        Direction::Right => vec![Direction::Up, Direction::Down],
    }
}

fn get_opposite_direction(direction: Direction) -> Direction {
    match direction {
        Direction::Up => Direction::Down,
        Direction::Down => Direction::Up,
        Direction::Left => Direction::Right,
        Direction::Right => Direction::Left,
    }
}

fn pathfinding(grid: &Grid, start: Point, goal: Point) -> Option<usize> {
    let dim = grid.len();
    let mut heap = BinaryHeap::new();
    let mut visited = vec![vec![false; grid.len()]; grid.len()];
    let mut cost_so_far = vec![vec![vec![vec![0;4];4]; grid.len()]; grid.len()];
    heap.push(State { point: start, origin: Origin { direction: Direction::Left, distance: 0 }, cost: 0});

    let mut counter = 0;

    while let Some(State { point, origin, cost }) = heap.pop() {
        counter += 1;
        if counter > 2 {
            break;
        }
        // Alternatively we could have continued to find all shortest paths
        if point == goal {
            return Some(cost);
        }

        // opposite direction
        let opposite_direction = get_opposite_direction(origin.direction);
        let new_point_opt = point.get_point_in_direction(opposite_direction, dim);
        let new_distance = origin.distance + 1;
        // Debug output oppsite direction
        println!("Opposite direction: {:?}, new_point_opt: {:?}, new_distance: {}", opposite_direction, new_point_opt, new_distance);
        if new_point_opt.is_some() && new_distance < 4 {
            let new_point = new_point_opt.unwrap();
            let new_cost = cost + grid[new_point.y][new_point.x] as usize;
            if new_cost < cost_so_far[new_point.y][new_point.x][origin.direction as usize][new_distance] {
                cost_so_far[new_point.y][new_point.x][origin.direction as usize][new_distance] = new_cost;
                heap.push(State { point: new_point, origin: Origin { direction: origin.direction, distance: new_distance }, cost: new_cost });
            }
        }

        // orthogonal directions
        for direction in get_orthogonal_directions(origin.direction) {
            let new_point_opt = point.get_point_in_direction(direction, dim);
            if new_point_opt.is_some() {
                let new_point = new_point_opt.unwrap();
                let new_cost = cost + grid[new_point.y][new_point.x] as usize;
                if new_cost < cost_so_far[new_point.y][new_point.x][get_opposite_direction(direction) as usize][1] {
                    cost_so_far[new_point.y][new_point.x][get_opposite_direction(direction) as usize][1] = new_cost;
                    heap.push(State { point: new_point, origin: Origin { direction: direction, distance: 0 }, cost: new_cost });
                }
            }
        }
    }
    // print heap
    println!("{:?}", heap);
    None
}

fn get_path(previous: &[Vec<Option<Point>>], start: Point, goal: Point) -> Vec<Point> {
    let mut path = Vec::new();
    let mut current = goal;
    while current != start {
        path.push(current);
        current = previous[current.y][current.x].unwrap();
    }
    path.push(start);
    path.reverse();
    path
}

fn print_grid(grid: &Grid, path: &[Point]) {
    for (i, row) in grid.iter().enumerate() {
        for (j, cell) in row.iter().enumerate() {
            if path.contains(&Point { x: j, y: i }) {
                print!("{} ", cell);
            } else {
                print!(". ");
            }
        }
        println!();
    }
}

fn solution_part_1() {
    let grid = load_data("input_example.txt");
    // grid[0][3] = 20;
    let goal = Point { x: grid.len() - 1, y: grid.len() - 1 };
    let start = Point { x: 0, y: 0 };
    let cost = pathfinding(&grid, start, goal);
    println!("Part 1: {}", cost.is_some());
}

fn main() {
    let start = std::time::Instant::now();
    solution_part_1();
    println!("Finished in {} us", start.elapsed().as_micros());
}
