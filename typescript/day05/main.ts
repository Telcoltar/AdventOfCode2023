import fs from "fs";

type OffsetInterval = Interval & {
  offset: number;
};

type Interval = {
  start?: number;
  end?: number;
};

function isInInterval(interval: Interval, number: number): boolean {
  if (interval.start && interval.end) {
    return number >= interval.start && number <= interval.end;
  }
  if (interval.start) {
    return number >= interval.start;
  }
  if (interval.end) {
    return number <= interval.end;
  }
  return false;
}

function intersectOffsetIntervalWithInterval(
  interval: Interval,
  offsetInterval: OffsetInterval,
): Interval {
  if (offsetInterval.start && offsetInterval.end) {
    return {
      start:
        Math.max(interval.start!, offsetInterval.start) + offsetInterval.offset,
      end: Math.min(interval.end!, offsetInterval.end) + offsetInterval.offset,
    };
  }
  if (offsetInterval.start) {
    return {
      start:
        Math.max(interval.start!, offsetInterval.start) + offsetInterval.offset,
      end: interval.end! + offsetInterval.offset,
    };
  }
  if (offsetInterval.end) {
    return {
      start: interval.start! + offsetInterval.offset,
      end: Math.min(interval.end!, offsetInterval.end) + offsetInterval.offset,
    };
  }
  return {
    start: interval.start! + offsetInterval.offset,
    end: interval.end! + offsetInterval.offset,
  };
}

function loadData(filename: string): [number[], OffsetInterval[][]] {
  const data = fs.readFileSync(filename).toString();
  const lines = data.split("\n");
  const seeds = parseSeeds(lines[0]);

  let currentLine = 2;
  const intervals: OffsetInterval[][] = [];
  while (currentLine < lines.length) {
    const { newLine, parsedData } = parseLines(lines, currentLine);
    currentLine = newLine + 1;
    intervals.push(parsedData);
  }
  return [seeds, intervals];
}

function fillGaps(intervals: OffsetInterval[]): OffsetInterval[] {
  intervals.sort((a, b) => a.start! - b.start!);
  let intervalsToAdd: OffsetInterval[] = [];
  if (intervals[0].start != 0) {
    intervalsToAdd.push({
      start: 0,
      end: intervals[0].start! - 1,
      offset: 0,
    });
  }
  for (let i = 0; i < intervals.length - 1; i++) {
    if (intervals[i].end! != intervals[i + 1].start! - 1) {
      intervalsToAdd.push({
        start: intervals[i].end! + 1,
        end: intervals[i + 1].start! - 1,
        offset: 0,
      });
    }
  }
  intervalsToAdd.push({
    start: intervals[intervals.length - 1].end! + 1,
    offset: 0,
  });
  intervals = intervals.concat(intervalsToAdd);
  intervals.sort((a, b) => a.start! - b.start!);
  return intervals;
}

// Parse the seeds line into number array
function parseSeeds(line: string): number[] {
  return line
    .split(":")[1]
    .trim()
    .split(" ")
    .map((x) => parseInt(x));
}

// Parse lines into OffsetInterval objects until an empty line is encountered
function parseLines(
  lines: string[],
  currentLine: number,
): { newLine: number; parsedData: OffsetInterval[] } {
  let parsedData: OffsetInterval[] = [];
  currentLine += 1;
  while (lines[currentLine] != "") {
    let line = lines[currentLine];
    let [dest, source, range] = line.split(" ").map((x) => parseInt(x));
    parsedData.push({
      start: source,
      end: source + range - 1,
      offset: dest - source,
    });
    currentLine += 1;
  }
  parsedData = fillGaps(parsedData);
  return { newLine: currentLine, parsedData };
}

function solutionPart1() {
  let [seeds, maps] = loadData("input.txt");
  let currentSeeds = seeds.toSorted((a: number, b: number) => a - b);
  for (let seedMap of maps) {
    let currentMapIndex = 0;
    let nextSeeds: number[] = [];
    for (let seed of currentSeeds) {
      while (!isInInterval(seedMap[currentMapIndex], seed)) {
        currentMapIndex += 1;
      }
      nextSeeds.push(seed + seedMap[currentMapIndex].offset);
    }
    nextSeeds.sort((a: number, b: number) => a - b);
    currentSeeds = nextSeeds;
  }
  console.log(currentSeeds[0]);
}

function solutionPart2() {
  let [seeds, maps] = loadData("input.txt");
  let currentSeedIntervals: Interval[] = [];
  for (let i = 0; i < seeds.length - 1; i += 2) {
    currentSeedIntervals.push({
      start: seeds[i],
      end: seeds[i] + seeds[i + 1] - 1,
    });
  }
  currentSeedIntervals.sort((a, b) => a.start! - b.start!);
  for (let seedMap of maps) {
    let currentMapIndex = 0;
    let nextSeedIntervals: Interval[] = [];
    for (let seedInterval of currentSeedIntervals) {
      while (!isInInterval(seedMap[currentMapIndex], seedInterval.start!)) {
        currentMapIndex += 1;
      }
      let start = currentMapIndex;
      while (!isInInterval(seedMap[currentMapIndex], seedInterval.end!)) {
        currentMapIndex += 1;
      }
      for (let i = start; i <= currentMapIndex; i++) {
        nextSeedIntervals.push(
          intersectOffsetIntervalWithInterval(seedInterval, seedMap[i]),
        );
      }
    }
    nextSeedIntervals.sort((a, b) => a.start! - b.start!);
    currentSeedIntervals = nextSeedIntervals;
  }
  console.log(currentSeedIntervals[0].start);
}

solutionPart1();
solutionPart2();
