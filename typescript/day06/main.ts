import fs from 'fs';

function loadDataPart1(filename: string): [number[], number[]] {
  const data = fs.readFileSync(filename, 'utf8').split("\n");
  const times = data[0].split(":")[1].trim().split(/\s+/).map((x) => parseInt(x));
  const distances = data[1].split(":")[1].trim().split(/\s+/).map((x) => parseInt(x));
  return [times, distances]
}

function loadDataPart2(filename: string): [number, number] {
  const data = fs.readFileSync(filename, 'utf8').split("\n");
  const time = parseInt(data[0].split(":")[1].trim().replaceAll(" ", ""))
  const distance = parseInt(data[1].split(":")[1].trim().replaceAll(" ", ""))
  return [time, distance]
}

function calculateRange(time: number, distance: number): number {
  const p_2 = time / 2
  const d = Math.sqrt(Math.pow(time / 2, 2) - distance)
  const zero_1 = p_2 - d
  const zero_2 = p_2 + d
  let high = Math.floor(zero_2)
  let low = Math.ceil(zero_1)
  if (zero_2 == high) {
    high -= 1
  }
  if (zero_1 == low) {
    low += 1
  }
  return high - low + 1
}
function solutionPart1() {
  const [times, distances] = loadDataPart1("input.txt");
  let ranges = []
  for (let i = 0; i < times.length; i++) {
    ranges.push(calculateRange(times[i], distances[i]))
  }
  console.log(ranges.reduce((a, b) => a * b, 1))
}

function solutionPart2() {
  const [time, distance] = loadDataPart2("input.txt");
  console.log(calculateRange(time, distance))
}

solutionPart1();
solutionPart2();