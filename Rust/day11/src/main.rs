
#[derive(Debug)]
struct Point {
    x: usize,
    y: usize,
}

impl Point {
    fn distance(&self, other: &Point) -> usize {
        let x = self.x.abs_diff(other.x);
        let y = self.y.abs_diff(other.y);
        x + y
    }

}

fn load_data(filename: &str) -> Vec<Point> {
    let data = std::fs::read_to_string(filename).unwrap();
    let mut galaxies = Vec::new();
    for (y, line) in data.lines().enumerate() {
        for (x, c) in line.chars().enumerate() {
            if c == '#' {
                galaxies.push(Point { x, y });
            }
        }
    }
    galaxies
}

macro_rules! make_new_galaxy_positions {
    ($func_name:ident, $coord:ident) => {
        fn $func_name(galaxies: &mut [Point], spread_factor: usize) {
            galaxies.sort_by_key(|g| g.$coord);
            let mut current_spread = 0;
            let mut last_pos = 0;
            for galaxy in galaxies {
                if (galaxy.$coord - last_pos) > 1 {
                    current_spread += ((galaxy.$coord - last_pos) - 1) * (spread_factor - 1);
                }
                last_pos = galaxy.$coord;
                galaxy.$coord += current_spread;
            }
        }
    };
}

make_new_galaxy_positions!(calculate_new_galaxy_positions_x, x);
make_new_galaxy_positions!(calculate_new_galaxy_positions_y, y);

fn calculate_new_galaxy_positions(galaxies: &mut [Point], spread_factor: usize) {
    calculate_new_galaxy_positions_y(galaxies, spread_factor);
    calculate_new_galaxy_positions_x(galaxies, spread_factor);
}

#[allow(dead_code)]
fn print_map(map: &[Vec<bool>]) {
    for row in map {
        for column in row {
            if *column {
                print!("#");
            } else {
                print!(".");
            }
        }
        println!();
    }
}


fn sum_distances(galaxy: &Point, galaxies: &[Point]) -> usize {
    galaxies.iter().map(|g| galaxy.distance(g)).sum()
}

fn solution(spread_factor: usize) {
    let mut galaxies = load_data("input.txt");
    calculate_new_galaxy_positions(&mut galaxies, spread_factor);
    let total_sum = galaxies.iter().map(|g| sum_distances(g, &galaxies)).sum::<usize>();
    println!("Total sum: {}", total_sum / 2);
}

fn solution_part_1() {
    solution(2);
}

fn solution_part_2() {
    solution(1000000);
}

fn main() {
    solution_part_1();
    solution_part_2();
}
