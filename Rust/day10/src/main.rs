use std::collections::HashSet;
use std::ops::{Index, IndexMut};

#[derive(Debug, Copy, Clone, Eq, PartialEq, Ord, PartialOrd, Hash)]
enum Direction {
    Up,
    Down,
    Left,
    Right,
}

impl Direction{
    fn opposite(&self) -> Direction {
        match self {
            Direction::Up => Direction::Down,
            Direction::Down => Direction::Up,
            Direction::Left => Direction::Right,
            Direction::Right => Direction::Left
        }
    }
}

fn get_symbol_for_direction(directions: (Direction, Direction)) -> char {
    println!("Directions: {:?}", directions);
    match directions {
        (Direction::Up, Direction::Down) => '|',
        (Direction::Left, Direction::Right) => '-',
        (Direction::Down, Direction::Right) => 'F',
        (Direction::Down, Direction::Left) => '7',
        (Direction::Up, Direction::Left) => 'J',
        (Direction::Up, Direction::Right) => 'L',
        _ => panic!("Unknown directions: {:?}", directions)
    }
}

#[derive(Debug, Copy, Clone, Eq, PartialEq)]
enum Tile {
    Ground,
    Pipe(Pipe),
    Start
}

#[derive(Debug, Copy, Clone, Eq, PartialEq)]
struct Pipe {
    symbol: char,
    connections: (Direction, Direction)
}

impl Pipe {
    fn is_connected(&self, direction: Direction) -> bool {
        self.connections.0 == direction || self.connections.1 == direction
    }

    fn get_connected_direction(&self, direction: Direction) -> Direction {
        if self.connections.0 == direction {
            self.connections.1
        } else if self.connections.1 == direction {
            self.connections.0
        } else {
            panic!("No connection in direction {:?}", direction);
        }
    }
}

#[derive(Debug, Copy, Clone, Eq, PartialEq, Hash)]
struct Point {
    x: usize,
    y: usize
}

impl Point {
    fn above(&self) -> Point {
        Point { x: self.x, y: self.y - 1 }
    }

    fn below(&self) -> Point {
        Point { x: self.x, y: self.y + 1 }
    }

    fn left(&self) -> Point {
        Point { x: self.x - 1, y: self.y }
    }

    fn right(&self) -> Point {
        Point { x: self.x + 1, y: self.y }
    }

    fn neighbour_in_direction(&self, direction: Direction) -> Point {
        match direction {
            Direction::Up => self.above(),
            Direction::Down => self.below(),
            Direction::Left => self.left(),
            Direction::Right => self.right()
        }
    }
}

type Map = Vec<Vec<Tile>>;

impl Index<Point> for Map {
    type Output = Tile;

    fn index(&self, index: Point) -> &Self::Output {
        &self[index.y][index.x]
    }
}

impl Index<&Point> for Map {
    type Output = Tile;

    fn index(&self, index: &Point) -> &Self::Output {
        &self[index.y][index.x]
    }
}

impl IndexMut<Point> for Map {
    fn index_mut(&mut self, index: Point) -> &mut Self::Output {
        &mut self[index.y][index.x]
    }
}
fn load_data(filename: &str) ->  Map {
    let data = std::fs::read_to_string(filename).unwrap();
    let mut tiles = Vec::new();
    for line in data.lines() {
        let mut line_tiles = Vec::new();
        for c in line.chars() {
            match c {
                '.' => line_tiles.push(Tile::Ground),
                '|' => line_tiles.push(Tile::Pipe(Pipe { symbol: '|', connections: (Direction::Up, Direction::Down) })),
                '-' => line_tiles.push(Tile::Pipe(Pipe { symbol: '-', connections: (Direction::Left, Direction::Right) })),
                'F' => line_tiles.push(Tile::Pipe(Pipe { symbol: 'F', connections: (Direction::Right, Direction::Down) })),
                '7' => line_tiles.push(Tile::Pipe(Pipe { symbol: '7', connections: (Direction::Down, Direction::Left) })),
                'J' => line_tiles.push(Tile::Pipe(Pipe { symbol: 'J', connections: (Direction::Left, Direction::Up) })),
                'L' => line_tiles.push(Tile::Pipe(Pipe { symbol: 'L', connections: (Direction::Up, Direction::Right) })),
                'S' => line_tiles.push(Tile::Start),
                _ => panic!("Unknown tile type: {}", c)
            }
        }
        tiles.push(line_tiles);
    }
    pad_tiles(&mut tiles);
    tiles
}

fn pad_tiles(tiles: &mut Map) {
    let width = tiles[0].len();
    for line in tiles.iter_mut() {
        line.insert(0, Tile::Ground);
        line.push(Tile::Ground);
    }
    let mut ground_line = Vec::new();
    for _ in 0..width + 2 {
        ground_line.push(Tile::Ground);
    }
    tiles.insert(0, ground_line.clone());
    tiles.push(ground_line);
}

fn find_start(tiles: &[Vec<Tile>]) -> Point {
    for (y, line) in tiles.iter().enumerate() {
        for (x, tile) in line.iter().enumerate() {
            if let Tile::Start = tile { return Point{x, y } }
        }
    }
    panic!("No start found");
}

fn get_connected_directions(tiles: &Map, point: &Point) -> Vec<Direction> {
    let mut directions = Vec::new();
    for direction in [Direction::Up, Direction::Down, Direction::Left, Direction::Right] {
        if let Tile::Pipe(pipe) = tiles[point.neighbour_in_direction(direction)] {
            if pipe.is_connected(direction.opposite()) {
                directions.push(direction);
            }
        }
    }
    directions
}

fn follow_path(start: Point, direction: Direction, end: Point, tiles: &Map) -> HashSet<Point> {
    let mut path = HashSet::new();
    path.insert(start);
    let mut current_direction = direction;
    let mut current = start;
    loop {
        current = current.neighbour_in_direction(current_direction);
        path.insert(current);
        if let Tile::Pipe(pipe) = tiles[current] {
            current_direction = pipe.get_connected_direction(current_direction.opposite());
            if current == end {
                break;
            }
        } else {
            break;
        }
    }
    path
}
fn solution_part_1() {
    let data = load_data("input.txt");
    let start = find_start(&data);
    println!("Start: {:?}", start);
    let directions = get_connected_directions(&data, &start);
    println!("Directions: {:?}", directions);
    for direction in directions {
        let path = follow_path(start, direction, start, &data);
        println!("Path len: {:?}", path.len());
        println!("Middle: {:?}", path.len() / 2);
    }
}

fn scan_line(y: usize, line: &[Tile], path: &HashSet<Point>) -> Vec<Point> {
    let mut points = Vec::new();
    let mut is_inside = false;
    let mut i = 1;
    while i < line.len() - 1 {
        let tile = line[i];
        // println!("Tile: {:?}, y: {:?}, x: {:?}, inside: {:?}, path: {:?}", tile, y, i, is_inside, path.contains(&Point { x: i, y }));
        if let Tile::Pipe(pipe) = tile {
            if !path.contains(&Point { x: i, y }) {
                if is_inside {
                    points.push(Point { x: i, y });
                }
            } else if pipe.symbol == '|' {
                is_inside = !is_inside;
            } else if pipe.symbol == 'F' {
                i += 1;
                while let Tile::Pipe(pipe) = line[i] {
                    if pipe.symbol == '-' {
                        i += 1;
                    } else {
                        break;
                    }
                }
                let tile = line[i];
                if let Tile::Pipe(pipe) = tile {
                    if pipe.symbol == 'J' {
                        is_inside = !is_inside;
                    }
                }
            } else if pipe.symbol == 'L' {
                i += 1;
                while let Tile::Pipe(pipe) = line[i] {
                    if pipe.symbol == '-' {
                        i += 1;
                    } else {
                        break;
                    }
                }
                let tile = line[i];
                if let Tile::Pipe(pipe) = tile {
                    if pipe.symbol == '7' {
                        is_inside = !is_inside;
                    }
                }
            }
        } else if is_inside {
            points.push(Point { x: i, y });
        }
        i += 1;
    }
    points
}

fn scan_column(x: usize, column: &[Tile], path: &HashSet<Point>) -> Vec<Point> {
    let mut points = Vec::new();
    let mut is_inside = false;
    let mut i = 1;
    while i < column.len() - 1 {
        let tile = column[i];
        // println!("Tile: {:?}, y: {:?}, x: {:?}, inside: {:?}", tile, i, x, is_inside);
        if let Tile::Pipe(pipe) = tile {
            if !path.contains(&Point { x, y: i }) {
                if is_inside {
                    points.push(Point { x, y: i });
                }
            } else if pipe.symbol == '-' {
                is_inside = !is_inside;
            } else if pipe.symbol == 'F' {
                i += 1;
                while let Tile::Pipe(pipe) = column[i] {
                    if pipe.symbol == '|' {
                        i += 1;
                    } else {
                        break;
                    }
                }
                let tile = column[i];
                if let Tile::Pipe(pipe) = tile {
                    if pipe.symbol == 'J' {
                        is_inside = !is_inside;
                    }
                }
            } else if pipe.symbol == '7' {
                i += 1;
                while let Tile::Pipe(pipe) = column[i] {
                    if pipe.symbol == '|' {
                        i += 1;
                    } else {
                        break;
                    }
                }
                let tile = column[i];
                if let Tile::Pipe(pipe) = tile {
                    if pipe.symbol == 'L' {
                        is_inside = !is_inside;
                    }
                }
            }
        } else if is_inside {
            points.push(Point { x, y: i });
        }
        i += 1;
    }
    points
}

fn replace_start(tiles: &mut Map, start: Point) {
    let mut directions = get_connected_directions(tiles, &start);
    directions.sort();
    let symbol = get_symbol_for_direction((directions[0], directions[1]));
    tiles[start] = Tile::Pipe(Pipe { symbol, connections: (directions[0], directions[1]) });
}

fn flip_map(tiles: &mut Map) {
    let mut new_tiles = Vec::new();
    for x in 0..tiles[0].len() {
        let mut column = Vec::new();
        for y in 0..tiles.len() {
            column.push(tiles[y][x]);
        }
        new_tiles.push(column);
    }
    *tiles = new_tiles;
}

fn solution_part_2() {
    let mut data = load_data("input.txt");
    let start = find_start(&data);
    replace_start(&mut data, start);
    let directions = get_connected_directions(&data, &start);
    let path = follow_path(start, directions[0], start, &data);
    let mut sum = 0;
    for (y, line) in data.iter().enumerate() {
        let points = scan_line(y, line, &path).len();
        sum += points;
        // println!("Line: {:?} Points: {:?}", y, points);
    }
    println!("Sum: {:?}", sum);
    flip_map(&mut data);
    let mut sum = 0;
    for (x, column) in data.iter().enumerate() {
        let points = scan_column(x, column, &path).len();
        sum += points;
        // println!("Column: {:?} Points: {:?}", x, points);
    }
    println!("Sum: {:?}", sum);
}
fn main() {
    solution_part_1();
    solution_part_2();
}
