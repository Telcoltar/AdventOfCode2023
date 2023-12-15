class Interval:
    def __init__(self, start: int | None, end: int | None) -> None:
        if start is None and end is None:
            raise TypeError("Interval requires at least start or end parameter")
        self.start = start
        self.end = end

    @classmethod
    def from_offset_interval(cls, offset_interval: "OffsetInterval") -> "Interval":
        return cls(start=offset_interval.start + offset_interval.offset,
                   end=offset_interval.end + offset_interval.offset)

    def __lt__(self, other):
        if self.start is not None and other.start is not None:
            return self.start < other.start
        elif self.end is not None and other.end is not None:
            return self.end < other.end
        raise TypeError("If half intervals are used, both need to have same end specified")

    def __eq__(self, other):
        if self.start is not None and other.start is not None:
            return self.start == other.start
        if self.end is not None and other.end is not None:
            return self.end == other.end
        raise TypeError("If half intervals are used, both need to have same end specified")

    def __contains__(self, item: int):
        if type(item) is not int:
            raise TypeError(f"item {item} is not an integer")
        if self.start is not None and self.end is not None:
            return self.start <= item <= self.end
        if self.start is not None:
            return self.start <= item
        if self.end is not None:
            return item <= self.end

    def __repr__(self):
        return f"Interval({self.start}, {self.end})"

    def __str__(self):
        return f"I({self.start}, {self.end})"

    def __add__(self, other: "Interval | int") -> "Interval":
        if type(other) is Interval:
            return Interval(self.start + other.start, self.end + other.end)
        if type(other) is int:
            return Interval(self.start + other, self.end + other)
        raise TypeError(f"Other needs to be of type Interval or int, got {type(other)}")

    def intersection(self, other: "Interval | OffsetInterval") -> "Interval":
        if type(other) is Interval:
            return Interval(max(self.start, other.start), min(self.end, other.end))
        if type(other) is OffsetInterval:
            return Interval(max(self.start, other.start), min(self.end, other.end)) + other.offset
        raise TypeError(f"Other needs to be of type Interval or OffsetInterval, got {type(other)}")


class OffsetInterval(Interval):

    def __init__(self, start: int | None, end: int | None, offset: int):
        super().__init__(start, end)
        self.offset = offset

    def __repr__(self):
        return f"Interval({self.start}, {self.end}, {self.offset})"

    def __str__(self):
        return f"I({self.start}, {self.end}, {self.offset})"


def add_zero_offset_intervalls(intervals: list[OffsetInterval]) -> list[OffsetInterval]:
    intervals_to_add = []
    if intervals[0].start != 0:
        intervals_to_add.append(OffsetInterval(0, intervals[0].start - 1, 0))
    for i in range(0, len(intervals) - 1):
        if intervals[i].end + 1 != intervals[i + 1].start:
            intervals_to_add.append(OffsetInterval(intervals[i].end + 1, intervals[i + 1].start - 1, 0))
    intervals_to_add.append(OffsetInterval(intervals[-1].end + 1, None, 0))
    intervals.extend(intervals_to_add)
    intervals.sort()
    return intervals


def load_data(filename: str) -> tuple[list[int], list[list[OffsetInterval]]]:
    maps: list[list[OffsetInterval]] = []
    with open(filename) as f:
        seeds = list(map(int, f.readline().strip().split(":")[1].strip().split()))
        f.readline()
        while f.readline():
            current_map: list[OffsetInterval] = []
            while True:
                line = f.readline().strip()
                if line == "":
                    break
                numbers = list(map(int, line.split()))
                interval = OffsetInterval(numbers[1], numbers[1] + numbers[2] - 1, numbers[0] - numbers[1])
                current_map.append(interval)
            current_map.sort()
            add_zero_offset_intervalls(current_map)
            maps.append(current_map)
    return seeds, maps


def solution_part_1():
    seeds, maps = load_data('input.txt')
    seeds.sort()
    current_seeds = seeds
    for seed_map in maps:
        map_index = 0
        seed_index = 0
        next_seeds: list[int] = []
        while seed_index < len(current_seeds):
            while current_seeds[seed_index] not in seed_map[map_index]:
                map_index += 1
                if map_index == len(seed_map):
                    break
            if map_index == len(seed_map):
                for i in range(seed_index, len(current_seeds)):
                    next_seeds.append(current_seeds[i])
                break
            next_seeds.append(current_seeds[seed_index] + seed_map[map_index].offset)
            seed_index += 1
        next_seeds.sort()
        current_seeds = next_seeds
    print(current_seeds[0])


def build_seed_intervals(seeds: list[int]) -> list[Interval]:
    intervals: list[Interval] = []
    for i in range(1, len(seeds), 2):
        intervals.append(Interval(seeds[i - 1], seeds[i - 1] + seeds[i] - 1))
    return intervals


def solution_part_2():
    seeds, maps = load_data('input.txt')
    current_seed_intervals: list[Interval] = build_seed_intervals(seeds)
    current_seed_intervals.sort()
    for seed_map in maps:
        map_index = 0
        seed_index = 0
        next_seed_intervals: list[Interval] = []
        while seed_index < len(current_seed_intervals):
            while current_seed_intervals[seed_index].start not in seed_map[map_index]:
                map_index += 1
            start_index = map_index
            while current_seed_intervals[seed_index].end not in seed_map[map_index]:
                map_index += 1
            for i in range(start_index, map_index + 1):
                next_seed_intervals.append(current_seed_intervals[seed_index].intersection(seed_map[i]))
            seed_index += 1
        next_seed_intervals.sort()
        current_seed_intervals = next_seed_intervals
    print(current_seed_intervals[0].start)


def main():
    solution_part_1()
    solution_part_2()


if __name__ == '__main__':
    main()
