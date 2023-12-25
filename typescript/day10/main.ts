import * as fs from "fs";

enum Direction {
    Up,
    Down,
    Left,
    Right,
}

const Directions = [Direction.Up, Direction.Down, Direction.Left, Direction.Right];

function opposite(direction: Direction): Direction {
    switch (direction) {
        case Direction.Up:
            return Direction.Down;
        case Direction.Down:
            return Direction.Up;
        case Direction.Left:
            return Direction.Right;
        case Direction.Right:
            return Direction.Left;
    }
}

class Point {
    x: number
    y: number

    constructor(x: number, y: number) {
        this.x = x;
        this.y = y;
    }

    equals(other: Point): boolean {
        return this.x == other.x && this.y == other.y;
    }

    string(): string {
        return `${this.x}_${this.y}`;
    }

    get_neighbour_in_direction(dir: Direction): Point {
        switch (dir) {
            case Direction.Up:
                return new Point(this.x, this.y - 1);
            case Direction.Down:
                return new Point(this.x, this.y + 1);
            case Direction.Left:
                return new Point(this.x - 1, this.y);
            case Direction.Right:
                return new Point(this.x + 1, this.y);
        }
    }
}

const SYMBOL_DIRECTIONS: Record<string, [Direction, Direction]> = {
    "|": [Direction.Up, Direction.Down],
    "-": [Direction.Left, Direction.Right],
    "L": [Direction.Right, Direction.Up],
    "J": [Direction.Left, Direction.Up],
    "7": [Direction.Left, Direction.Down],
    "F": [Direction.Right, Direction.Down]
}

const CODE_SYMBOLS: Record<string, string> = {
    "1100": "|",
    "0011": "-",
    "1001": "L",
    "1010": "J",
    "0110": "7",
    "0101": "F"
}

class Tile {
    symbol: string
    connections: Map<Direction, boolean>
    path: Map<Direction, Direction>

    constructor(symbol: string) {
        this.symbol = symbol;
        this.connections = new Map<Direction, boolean>();
        this.path = new Map<Direction, Direction>();
        if (symbol in SYMBOL_DIRECTIONS) {
            for (let dir of SYMBOL_DIRECTIONS[symbol]) {
                this.connections.set(dir, true);
            }
            this.path.set(SYMBOL_DIRECTIONS[symbol][0], SYMBOL_DIRECTIONS[symbol][1]);
            this.path.set(SYMBOL_DIRECTIONS[symbol][1], SYMBOL_DIRECTIONS[symbol][0]);
        }
    }

}

function pad_tiles(tile_map: Tile[][]): Tile[][] {
    let new_tile_map: Tile[][] = [];
    let new_tiles: Tile[] = [];
    for (let x = 0; x < tile_map[0].length + 2; x++) {
        new_tiles.push(new Tile("."));
    }
    new_tile_map.push(new_tiles);
    for (let y = 0; y < tile_map.length; y++) {
        new_tiles = [];
        new_tiles.push(new Tile("."));
        for (let x = 0; x < tile_map[y].length; x++) {
            new_tiles.push(tile_map[y][x]);
        }
        new_tiles.push(new Tile("."));
        new_tile_map.push(new_tiles);
    }
    new_tiles = [];
    for (let x = 0; x < tile_map[0].length + 2; x++) {
        new_tiles.push(new Tile("."));
    }
    new_tile_map.push(new_tiles);
    return new_tile_map;

}


function load_data(filename: string): Tile[][] {
    let tile_map: Tile[][] = [];
    let lines = fs.readFileSync(filename, "utf-8").split("\n");
    for (let line of lines) {
        let tiles: Tile[] = [];
        for (let c of line) {
            let tile = new Tile(c);
            tiles.push(tile);
        }
        tile_map.push(tiles);
    }
    tile_map = pad_tiles(tile_map);
    return tile_map;
}

function find_starting_point(tile_map: Tile[][]): Point {
    for (let y = 0; y < tile_map.length; y++) {
        for (let x = 0; x < tile_map[y].length; x++) {
            if (tile_map[y][x].symbol == "S") {
                return new Point(x, y);
            }
        }
    }
    throw "No starting point found";
}

function get_connected_directions(point: Point, tile_map: Tile[][]): string {
    let code = ["0", "0", "0", "0"];
    for (let dir of Directions) {
        let neighbour = point.get_neighbour_in_direction(dir);
        if (tile_map[neighbour.y][neighbour.x].connections.get(opposite(dir))) {
            code[dir] = "1";
        }
    }
    return code.join("");
}
function replace_starting_point(tile_map: Tile[][], point: Point) {
    tile_map[point.y][point.x] = new Tile(CODE_SYMBOLS[get_connected_directions(point, tile_map)]);
}

function follow_path(tile_map: Tile[][], point: Point): Set<string> {
    let current_point = point
    let current_direction = tile_map[point.y][point.x].connections.keys().next().value;
    let path: Set<string> = new Set<string>();
    path.add(current_point.string());
    while (true) {
        current_point = current_point.get_neighbour_in_direction(current_direction);
        current_direction = tile_map[current_point.y][current_point.x].path.get(opposite(current_direction));
        path.add(current_point.string());
        if (current_point.equals(point)) {
            break;
        }
    }
    return path
}

function scan_rows(tile_map: Tile[][], path: Set<string>): number {
    let inside_count = 0;
    for (let y = 1; y < tile_map.length - 1; y++) {
        let x = 1
        let is_inside = false;
        while (x < tile_map[y].length - 1) {
            if (!path.has(new Point(x, y).string())) {
                if (is_inside) {
                    inside_count += 1;
                }
            } else if (tile_map[y][x].symbol == "|") {
                is_inside = !is_inside;
            } else if (tile_map[y][x].symbol == "F") {
                x++;
                while (tile_map[y][x].symbol == "-") {
                    x++;
                }
                if (tile_map[y][x].symbol == "J") {
                    is_inside = !is_inside;
                }
            } else if (tile_map[y][x].symbol == "L") {
                x++;
                while (tile_map[y][x].symbol == "-") {
                    x++;
                }
                if (tile_map[y][x].symbol == "7") {
                    is_inside = !is_inside;
                }
            }
            x++;
        }
    }
    return inside_count
}

function solution_part_1() {
    let tile_map = load_data("input.txt");
    let starting_point = find_starting_point(tile_map);
    replace_starting_point(tile_map, starting_point);
    let path = follow_path(tile_map, starting_point);
    console.log(path.size / 2);
}

function solution_part_2() {
    let tile_map = load_data("input.txt");
    let starting_point = find_starting_point(tile_map);
    replace_starting_point(tile_map, starting_point);
    let path = follow_path(tile_map, starting_point);
    let inside_count = scan_rows(tile_map, path);
    console.log(inside_count);
}

solution_part_1()
solution_part_2()