import * as fs from "fs";

function loadData(filename: string): number[][] {
    let lines = fs.readFileSync(filename, "utf-8").split("\n");
    return lines.map(line => line.split(" ").map(n => parseInt(n)));
}

function processLine(line: number[]): number[][] {
    let result: number[][] = [line];
    while (true) {
        let currentLine: number[] = []
        let lastLine = result[result.length - 1];
        for (let i = 0; i < lastLine.length - 1; i++) {
            currentLine.push(lastLine[i + 1] - lastLine[i]);
        }
        result.push(currentLine);
        if (currentLine.reduce((a, b) => Math.abs(a) + Math.abs(b), 0) == 0) {
            break;
        }
    }
    return result;
}

function predictEnd(lines: number[][]): number {
    return lines.map(line => line[line.length - 1]).reduce((a, b) => a + b, 0);
}

function predictStart(lines: number[][]): number {
    return lines.reverse().map(line => line[0]).reduce((acc, value) => value - acc, 0);
}

function solutionPart1() {
    let data = loadData("input.txt");
    console.log(data.map(line => predictEnd(processLine(line))).reduce((a, b) => a + b, 0));
}

function solutionPart2() {
    let data = loadData("input.txt");
    console.log(data.map(line => predictStart(processLine(line))).reduce((a, b) => a + b, 0));
}

solutionPart1()
solutionPart2()