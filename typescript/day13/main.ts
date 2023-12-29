import * as fs from "fs";


function load_data(filename: string): string[][] {
    return fs.readFileSync(filename, "utf8").split("\n\n").map(x => x.split("\n"));
}

function isStrictEqual(a: string, b: string): boolean {
    return  a === b;
}

function isSimilar(a: string, b: string): boolean {
    let diff = 0;
    for (let i = 0; i < a.length; i++) {
        if (a[i] != b[i]) {
            diff++;
        }
    }
    return diff == 1 || diff == 0
}

function get_start_points(grid: string[], equal_fn: (a: string, b: string) => boolean): number[] {
    let start_points: number[] = [];
    for (let i = 0; i < grid.length - 1; i++) {
        if (equal_fn(grid[i], grid[i+1])) {
            start_points.push(i);
        }
    }
    return start_points;
}

function check_starting_point(grid: string[], starting_point: number, equal_fn: (a: string, b: string) => boolean): boolean {
    let lower = starting_point
    let upper = starting_point + 1;
    while (lower >= 0 && upper < grid.length && equal_fn(grid[lower], grid[upper])) {
        lower--;
        upper++;
    }
    return lower == - 1 || upper == grid.length;
}

function scan_for_horizontal_reflection_line(grid: string[], equal_fn: (a: string, b: string) => boolean): number[] {
    let start_points = get_start_points(grid, equal_fn);
    let points: number[] = [];
    for (let i = 0; i < start_points.length; i++) {
        if (check_starting_point(grid, start_points[i], equal_fn)) {
            points.push(start_points[i]);
        }
    }
    return points
}

function transpose_grid(grid: string[]): string[] {
    let transposed_grid: string[] = [];
    for (let i = 0; i < grid[0].length; i++) {
        let row: string[] = [];
        for (let j = 0; j < grid.length; j++) {
            row.push(grid[j][i]);
        }
        transposed_grid.push(row.join(""));
    }
    return transposed_grid;
}

function solution(starting_point_fn: (grid: string[]) => number[]) {
    let grids = load_data("input.txt");
    let horizontal_sum = 0;
    let vertical_sum = 0;
    for (let grid of grids) {
        let horizontal_points = starting_point_fn(grid);
        if (horizontal_points.length == 1) {
            horizontal_sum += horizontal_points[0] + 1;
        } else {
            let transposed_grid = transpose_grid(grid);
            let vertical_points = starting_point_fn(transposed_grid);
            if (vertical_points.length == 1) {
                vertical_sum += vertical_points[0] + 1;
            }
        }
    }
    let total_sum = vertical_sum + horizontal_sum * 100
    console.log(total_sum)
}

function get_strict_starting_point(grid: string[]): number[] {
    return scan_for_horizontal_reflection_line(grid, isStrictEqual);
}

function get_loose_starting_point(grid: string[]): number[] {
    let strict_starting_points = scan_for_horizontal_reflection_line(grid, isStrictEqual);
    return scan_for_horizontal_reflection_line(grid, isSimilar).filter(x => !strict_starting_points.includes(x));
}

function solution_part_1() {
    solution(get_strict_starting_point);
}

function solution_part_2() {
    solution(get_loose_starting_point);
}

solution_part_1();
solution_part_2();