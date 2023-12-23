import * as fs from "fs";

enum Direction {
  L, R
}

function gcd(a: number, b: number): number {
    if (!b) return a;
    return gcd(b, a % b);
}
function load_data(filename: string): [Direction[], Map<string, Map<Direction, string>>] {
  let lines = fs.readFileSync(filename, 'utf8').split('\n');
  let directions = lines[0].split("").map((x) => x == "L" ? Direction.L : Direction.R);
  let network = new Map<string, Map<Direction, string>>();
  let matcher = /(\w{3}) = \((\w{3}), (\w+)\)/;
  for (let i = 2; i < lines.length; i++) {
    let match = matcher.exec(lines[i]);
    if (match) {
      let [_, source, left, right] = match;
      network.set(source, new Map<Direction, string>([[Direction.L, left], [Direction.R, right]]));
    }
  }
  return [directions, network];
}

function follow_path(start: string, directions: Direction[], network: Map<string, Map<Direction, string>>, check_fn: (node: string) => boolean): number {
  let current = start;
  let step = 0;
  let direction_index = 0;
  while (check_fn(current)) {
    current = network.get(current)!.get(directions[direction_index])!;
    step++;
    direction_index = (direction_index + 1) % directions.length;
  }
  return step;
}

function solution_part_1() {
  let [directions, network] = load_data("input.txt");
  let start = "AAA";
  let steps = follow_path(start, directions, network, (node) => node != "ZZZ");
  console.log(`Steps: ${steps}`);
}

function solution_part_2() {
  let [directions, network] = load_data("input.txt");
  let cycles = []
  for (let [node, _] of network) {
    if (node.endsWith("A")) {
      cycles.push(follow_path(node, directions, network, (node) => !node.endsWith("Z")));
    }
  }
  let lcm = cycles.reduce((a, b) => a * b / gcd(a, b));
  console.log(`Steps: ${lcm}`);
}

solution_part_1();
solution_part_2();