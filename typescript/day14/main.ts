import * as fs from "fs";

function load_data(file_name: string): string[][] {
    return fs.readFileSync(file_name, 'utf8').split('\n').map(line => line.split(''));
}

function tilt_north(grid: string[][]) {
    for (let row = 0; row < grid.length; row++) {
        for (let col = 0; col < grid[row].length; col++) {
            if (grid[row][col] == "O") {
                let current_row = row
                while (current_row > 0 && grid[current_row - 1][col] == ".") {
                    current_row--;
                }
                if (current_row != row) {
                    grid[current_row][col] = "O";
                    grid[row][col] = ".";
                }
            }
        }
    }
}

function calculate_load(grid: string[][]): number {
    let load = 0;
    for (let row = 0; row < grid.length; row++) {
        for (let col = 0; col < grid[row].length; col++)
            if (grid[row][col] == "O") {
                load += (grid.length - row);
            }
    }
    return load;
}

function solution_part_1() {
    let grid = load_data("input.txt");
    tilt_north(grid);
    console.log(calculate_load(grid));
}

function rotate_clockwise(grid: string[][]) {
    for (let row = 0; row < grid.length; row++) {
        for (let col = row + 1; col < grid[row].length; col++) {
            let temp = grid[row][col];
            grid[row][col] = grid[col][row];
            grid[col][row] = temp;
        }
    }
    for (let row = 0; row < grid.length; row++) {
        grid[row].reverse();
    }
}

function cycle(grid: string[][]) {
    for (let i = 0; i < 4; i++) {
        tilt_north(grid);
        rotate_clockwise(grid);
    }
}

function find_cycle(grid: string[][]): [number, number] {
    let cache = new Map<string, number>();
    cache.set(JSON.stringify(grid), 0);
    let current_cycle = 0;
    while (true) {
        cycle(grid);
        if (cache.has(JSON.stringify(grid))) {
            return [current_cycle - cache.get(JSON.stringify(grid))!, current_cycle + 1];
        }
        cache.set(JSON.stringify(grid), current_cycle);
        current_cycle++;
    }
}

function solution_part_2() {
    let grid = load_data("input.txt");
    let [cycle_length, current_cycle] = find_cycle(grid);
    let remaining_cycles = (1000000000 - current_cycle) % cycle_length
    for (let i = 0; i < remaining_cycles; i++) {
        cycle(grid);
    }
    console.log(calculate_load(grid));
}

solution_part_1()
solution_part_2()