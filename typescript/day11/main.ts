import * as fs from "fs";

type Point = {
    x: number,
    y: number
}

const coordinates: (keyof Point)[] = ["x", "y"];

function load_data(filename: string): Point[] {
    let lines = fs.readFileSync(filename, 'utf-8').split("\n");
    let galaxies: Point[] = []
    lines.forEach((line, y) => {
        for (let x = 0; x < line.length; x++) {
            if (line[x] == '#') {
                galaxies.push({ x: x, y: y });
            }
        }
    });
    return galaxies;
}

function spread_galaxies(galaxies: Point[], spread_factor: number) {
    for (let coord of coordinates) {
        galaxies.sort((a, b) => a[coord] - b[coord]);
        let current_spread = 0;
        let current_pos = 0;
        for (let galaxy of galaxies) {
            let diff = galaxy[coord] - current_pos;
            current_pos = galaxy[coord];
            if (diff > 1) {
                current_spread += (diff - 1) * (spread_factor - 1);
            }
            galaxy[coord] += current_spread;
        }
    }
}

function sum_distance(galaxies: Point[]): number {
    let sum = 0;
    for (let i = 0; i < galaxies.length; i++) {
        for (let j = i + 1; j < galaxies.length; j++) {
            sum += Math.abs(galaxies[i].x - galaxies[j].x) + Math.abs(galaxies[i].y - galaxies[j].y);
        }
    }
    return sum;
}

function solution(filename: string, spread_factor: number): number {
    let galaxies = load_data(filename);
    spread_galaxies(galaxies, spread_factor);
    return sum_distance(galaxies);
}

function solution_part_1(filename: string) {
    console.log(solution(filename, 2));
}

function solution_part_2(filename: string) {
    console.log(solution(filename, 1000000));
}

solution_part_1("input.txt");
solution_part_2("input.txt");