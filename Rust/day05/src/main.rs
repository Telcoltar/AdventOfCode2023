use std::cmp::min;
use std::fs;


#[derive(Debug, Eq, PartialEq, Ord, PartialOrd)]
struct MapNumber {
    value: i64,
    point_type: String,
    offset: i64,
}

#[derive(Debug, Eq, PartialEq, Ord, PartialOrd)]
struct Range {
    start: i64,
    end: i64,
}

fn load_data(filename: &str) -> (Vec<i64>, Vec<Vec<MapNumber>>) {
    let data = fs::read_to_string(filename).expect("Unable to read file");
    let mut lines = data.lines();
    let seeds = lines.next().unwrap().split(":").nth(1).unwrap().trim();
    let seeds = seeds.split_whitespace().map(|x| x.parse::<i64>().unwrap()).collect::<Vec<i64>>();
    _ = lines.next();
    _ = lines.next();
    let mut maps: Vec<Vec<MapNumber>> = Vec::new();
    let mut current_map: Vec<MapNumber> = Vec::new();
    while let Some(line) = lines.next() {
        if line == "" {
            current_map.sort();
            maps.push(current_map);
            current_map = Vec::new();
            _ = lines.next();
            continue
        }
        let line_split = line.split_whitespace().map(|x| x.parse::<i64>().unwrap()).collect::<Vec<i64>>();
        let point = MapNumber{value: line_split[1], point_type: String::from("start"), offset: line_split[0] - line_split[1]};
        current_map.push(point);
        let point = MapNumber{value: line_split[1] + line_split[2] - 1, point_type: String::from("end"), offset: line_split[0] - line_split[1]};
        current_map.push(point);
    }
    current_map.sort();
    maps.push(current_map);
    return (seeds, maps)
}

fn process_seed(seed: i64, maps: &Vec<Vec<MapNumber>>) -> i64 {
    let mut current_seed = seed;
    for map in maps.iter() {
        let mut i = 1;
        while current_seed > map[i].value && i < map.len() - 1 {
            i += 2;
        }
        // println!("Number: {}, {}, {}", map[i].value, map[i].point_type, map[i].offset);
        if map[i - 1].value <= current_seed && current_seed <= map[i].value {
            current_seed += map[i].offset;
        }
        // println!("Seed: {}", current_seed)
    }
    return current_seed
}
fn solution_part_1() {
    let (seeds, maps) = load_data("input.txt");
    let mut locations: Vec<i64> = Vec::new();
    for seed in seeds {
        let current_seed = process_seed(seed, &maps);
        locations.push(current_seed);
    }
    locations.sort();
    println!("Location: {}", locations[0]);
}

fn solution_part_2() {
    let (seeds, maps) = load_data("input.txt");
    let mut seed_ranges: Vec<Range> = Vec::new();
    for i in (0..seeds.len()).step_by(2) {
        seed_ranges.push(Range{start: seeds[i], end: seeds[i + 1] + seeds[i] - 1});
    }
    seed_ranges.sort();
    for map in maps {
        let mut current_ranges: Vec<Range> = Vec::new();
        let mut current_index = 0;
        for range in seed_ranges {
            let mut current_start = range.start;
            // search next value in map greater or equal to current_start
            while current_index < map.len() && map[current_index].value < current_start {
                current_index += 1;
            }
            // if current_index is out of bounds, break
            if current_index >= map.len() {
                current_ranges.push(Range{start: current_start, end: range.end});
                continue
            }
            loop {
                if current_index >= map.len() {
                    current_ranges.push(Range{start: current_start, end: range.end});
                    break
                }
                if map[current_index].point_type == "start" {
                    if current_start == map[current_index].value {
                        let current_end = min(range.end, map[current_index + 1].value);
                        current_ranges.push(Range{start: current_start + map[current_index].offset, end: current_end + map[current_index].offset});
                        if current_end == range.end {
                            break
                        }
                        current_start = current_end + 1;
                        current_index += 2;
                    } else {
                        let current_end = min(range.end, map[current_index].value - 1);
                        current_ranges.push(Range{start: current_start, end: current_end});
                        if current_end == range.end {
                            break
                        }
                        current_start = current_end + 1;
                    }
                } else if map[current_index].point_type == "end" {
                    if current_start == map[current_index].value {
                        current_ranges.push(Range{start: current_start + map[current_index].offset, end: current_start + map[current_index].offset});
                        current_start += 1;
                        current_index += 1;
                    } else {
                        let current_end = min(range.end, map[current_index].value);
                        current_ranges.push(Range{start: current_start + map[current_index].offset, end: current_end + map[current_index].offset});
                        if current_end == range.end {
                            break
                        }
                        current_start = current_end + 1;
                        current_index += 1;
                    }
                }
            }
        }
        current_ranges.sort();
        seed_ranges = current_ranges;
    }
    println!("Lowest value: {}", seed_ranges[0].start);
}
fn main() {
    solution_part_1();
    solution_part_2()
}
