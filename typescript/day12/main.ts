import * as fs from "fs";

enum Condition {
    Operational = ".",
    Damaged = "#",
    Unknown = "?"
}

function loadData(filename: string): [Condition[], number[]][] {
    let data = fs.readFileSync(filename, "utf8").split("\n");
    let result: [Condition[], number[]][] = [];
    for (let line of data) {
        let [conditionsStr, numbersStr] = line.split(" ");
        let conditions: Condition[] = [];
        for (let char of conditionsStr) {
            switch (char) {
                case ".": conditions.push(Condition.Operational); break;
                case "#": conditions.push(Condition.Damaged); break;
                case "?": conditions.push(Condition.Unknown); break;
            }
        }
        let numbers: number[] = numbersStr.split(",").map(Number);
        result.push([conditions, numbers]);
    }
    return result;
}

class Process {
    conditions: Condition[];
    blocks: number[];
    cache: number[][];

    constructor(conditions: Condition[], blocks: number[], factor: number = 1) {
        if (factor > 1) {
            this.multiplyInput(conditions, blocks, factor);
        }
        this.conditions = conditions;
        this.blocks = blocks;
        this.cache = Array.from({length: this.conditions.length + 2})
            .map(() => Array.from({length: this.blocks.length + 2}).map(() => -1));
    }

    multiplyInput(conditions: Condition[], blocks: number[], factor: number) {
        let originalConditions = conditions.slice();
        let originalBlocks = blocks.slice();
        for (let i = 0; i < factor - 1; i++) {
            conditions.push(Condition.Unknown)
            conditions.push(...originalConditions);
            blocks.push(...originalBlocks);
        }
    }

    processRow(condIndex: number, blockIndex: number): number {
        if (this.cache[condIndex][blockIndex] != -1) {
            return this.cache[condIndex][blockIndex];
        }
        if (condIndex >= this.conditions.length) {
            if (blockIndex >= this.blocks.length) {
                return 1;
            }
            return 0;
        }
        let result = 0;
        switch (this.conditions[condIndex]) {
            case Condition.Operational:
                result = this.processRow(condIndex + 1, blockIndex);
                break
            case Condition.Damaged:
                if (blockIndex < this.blocks.length && this.hasBlockSpace(condIndex, this.blocks[blockIndex])) {
                    result = this.processRow(condIndex + this.blocks[blockIndex] + 1, blockIndex + 1);
                }
                break
            case Condition.Unknown:
                result = this.processRow(condIndex + 1, blockIndex);
                if (blockIndex < this.blocks.length && this.hasBlockSpace(condIndex, this.blocks[blockIndex])) {
                    result += this.processRow(condIndex + this.blocks[blockIndex] + 1, blockIndex + 1);
                }
        }
        this.cache[condIndex][blockIndex] = result;
        return result;
    }

    hasBlockSpace(condIndex: number, block: number) {
        while (condIndex < this.conditions.length && block > 0 &&
            (this.conditions[condIndex] == Condition.Unknown || this.conditions[condIndex] == Condition.Damaged)) {
            condIndex++;
            block--;
        }
        return block == 0 && (condIndex == this.conditions.length
            || this.conditions[condIndex] == Condition.Unknown || this.conditions[condIndex] == Condition.Operational);
    }

}

function solutionPart1() {
    solution()
}

function solutionPart2() {
    solution(5)
}

function solution(factor: number = 1): void {
    let data = loadData("input.txt");
    let result = 0;
    let start = new Date().getTime();
    for (let [conditions, numbers] of data) {
        let p = new Process(conditions, numbers, factor);
        result += p.processRow(0, 0);
    }
    let end = new Date().getTime();
    console.log(`Time: ${end - start} ms`);
    console.log(result);
}

solutionPart1()
solutionPart2()
