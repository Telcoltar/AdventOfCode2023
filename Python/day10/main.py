from dataclasses import dataclass
from enum import Enum, auto


class Direction(Enum):
    Up = auto()
    Down = auto()
    Left = auto()
    Right = auto()

    def get_opposite(self):
        if self == Direction.Up:
            return Direction.Down
        elif self == Direction.Down:
            return Direction.Up
        elif self == Direction.Left:
            return Direction.Right
        elif self == Direction.Right:
            return Direction.Left


@dataclass
class Point:
    x: int
    y: int

    def get_next_in_direction(self, direction: Direction):
        if direction == Direction.Up:
            return Point(self.x, self.y - 1)
        elif direction == Direction.Down:
            return Point(self.x, self.y + 1)
        elif direction == Direction.Left:
            return Point(self.x - 1, self.y)
        elif direction == Direction.Right:
            return Point(self.x + 1, self.y)

    def __hash__(self):
        return hash((self.x, self.y))


SYMBOL_DIRECTIONS = {
    "|": (Direction.Up, Direction.Down),
    "-": (Direction.Left, Direction.Right),
    "L": (Direction.Right, Direction.Up),
    "J": (Direction.Left, Direction.Up),
    "7": (Direction.Left, Direction.Down),
    "F": (Direction.Right, Direction.Down)
}


class Tile:

    def __init__(self, symbol: str):
        self.symbol = symbol
        self.connection_dict = {d: False for d in Direction}
        if symbol in SYMBOL_DIRECTIONS:
            self.tile_type = "pipe"
            self.connections = SYMBOL_DIRECTIONS[symbol]
            for conn in self.connections:
                self.connection_dict[conn] = True
            self.path_dict = {self.connections[0]: self.connections[1], self.connections[1]:self.connections[0]}
        elif symbol == ".":
            self.tile_type = "ground"
        elif symbol == "S":
            self.tile_type = "start"


    def __str__(self):
        return str([int(value) for value in self.connection_dict.values()])

    def __repr__(self):
        return self.__str__()



def pad_tiles(tiles: list[list[Tile]]) -> list[list[Tile]]:
    padded_tiles = []
    width = len(tiles[0])
    padded_tiles.append([Tile(".") for _ in range(width + 2)])
    for line in tiles:
        padded_tiles.append([Tile(".")] + line + [Tile(".")])
    padded_tiles.append([Tile(".") for _ in range(width + 2)])
    return padded_tiles


def get_symbol_from_direction(direction: dict[Direction, bool]) -> str:
    code = "".join([str(int(direction[d])) for d in Direction])
    if code == "1100":
        return "|"
    elif code == "0011":
        return "-"
    elif code == "1001":
        return "L"
    elif code == "1010":
        return "J"
    elif code == "0110":
        return "7"
    elif code == "0101":
        return "F"




class Map:

    def __init__(self, tiles: list[list[Tile]]) -> None:
        self.tiles = pad_tiles(tiles)
        self.start = self.find_start()
        directions = {d: False for d in Direction}
        for d in self.get_connected_directions(self.start):
            directions[d] = True
        self[self.start] = Tile(get_symbol_from_direction(directions))

    def get_connected_directions(self, point: Point) -> list[Direction]:
        directions = []
        for d in Direction:
            neighbour = self[point.get_next_in_direction(d)]
            if neighbour.connection_dict[d.get_opposite()]:
                directions.append(d)
        return directions

    def find_start(self) -> Point:
        for y, line in enumerate(self.tiles):
            for x, tile in enumerate(line):
                if tile.tile_type == "start":
                    return Point(x, y)

    def __getitem__(self, item):
        if isinstance(item, Point):
            return self.tiles[item.y][item.x]
        return self.tiles[item]

    def __setitem__(self, key, value):
        if isinstance(key, Point) and isinstance(value, Tile):
            self.tiles[key.y][key.x] = value
        else:
            self.tiles[key] = value

    def __repr__(self):
        return "\n".join(["".join([str(tile) for tile in line]) for line in self.tiles])

    def __len__(self):
        return len(self.tiles)


def load_data(filename: str) -> Map:
    tiles = []
    with open(filename) as f:
        for line in f.readlines():
            tiles.append([Tile(c) for c in line.strip()])
    return Map(tiles)


def follow_path(start: Point, end: Point, tile_map: Map) -> set[Point]:
    current = start
    direction = tile_map[start].connections[0]
    path: set[Point] = {current}
    while True:
        current = current.get_next_in_direction(direction)
        direction = tile_map[current].path_dict[direction.get_opposite()]
        path.add(current)
        if current == end:
            break
    return path


def scan_rows(tile_map: Map, path: set[Point]) -> int:
    inside_point: set[Point] = set()
    height = len(tile_map)
    width = len(tile_map[0])
    for y in range(1, height - 1):
        is_inside = False
        x = 1
        while x < width - 1:
            point = Point(x=x, y=y)
            if point not in path:
                if is_inside:
                    inside_point.add(point)
            elif tile_map[point].symbol == "|":
                is_inside = not is_inside
            elif tile_map[point].symbol == "F":
                x += 1
                while tile_map[Point(x,y)].symbol == "-":
                    x += 1
                if tile_map[Point(x, y)].symbol == "J":
                    is_inside = not is_inside
            elif tile_map[point].symbol == "L":
                x += 1
                while tile_map[Point(x,y)].symbol == "-":
                    x += 1
                if tile_map[Point(x, y)].symbol == "7":
                    is_inside = not is_inside
            x += 1
    return len(inside_point)


def solution_part_1():
    tile_map = load_data("input.txt")
    path = follow_path(tile_map.start, tile_map.start, tile_map)
    print(len(path) // 2)


def solution_part_2():
    tile_map = load_data("input.txt")
    path = follow_path(tile_map.start, tile_map.start, tile_map)
    result = scan_rows(tile_map, path)
    print(result)


def main():
    solution_part_1()
    solution_part_2()


if __name__ == "__main__":
    main()
