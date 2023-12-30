import * as fs from "fs";

function load_data(filename: string): string[] {
    return fs.readFileSync(filename, 'utf8').split(',');
}

function hash(input: string): number {
    let hash_sum = 0;
    for (let i = 0; i < input.length; i++) {
        hash_sum += input.charCodeAt(i);
        hash_sum *= 17;
        hash_sum %= 256;
    }
    return hash_sum
}

function solution_part_1() {
    let instructions = load_data('input.txt');
    let total_sum = instructions.map(hash).reduce((a, b) => a + b);
    console.log(total_sum);
}

type Box = {
    index_map: Map<string, number>
    values: number[]
    current_index: number
}

function create_box(): Box {
    return {
        index_map: new Map(),
        values: [],
        current_index: 0
    }
}

function add_value(box: Box, label: string, value: number) {
    if (box.index_map.has(label)) {
        let index = box.index_map.get(label)!;
        box.values[index] = value;
    } else {
        box.index_map.set(label, box.current_index);
        box.values.push(value);
        box.current_index += 1;
    }
}

function remove_value(box: Box, label: string) {
    if (box.index_map.has(label)) {
        let index = box.index_map.get(label)!;
        box.values[index] = -1;
        box.index_map.delete(label);
    }
}

function solution_part_2() {
    let instructions = load_data('input.txt');
    let boxes: Box[] = Array.from({length: 256}, create_box);
    let instruction_matcher = /([a-z]+)([=-])(\d*)/;
    for (let instruction of instructions) {
        let match = instruction_matcher.exec(instruction);
        if (match) {
            let label = match[1];
            let operation = match[2];
            let box = boxes[hash(label)];
            if (operation === '=') {
                let value = parseInt(match[3]);
                add_value(box, label, value);
            } else if (operation === '-') {
                remove_value(box, label);
            }
        }
    }
    let total_sum = 0;
    boxes.forEach((box, index) => {
       if (box.current_index > 0) {
           let box_sum = 0;
           let current_index = 1;
           for (let value of box.values) {
               if (value != -1) {
                   box_sum += value * current_index;
                   current_index += 1;
               }
           }
              total_sum += box_sum * (index + 1);
       }
    });
    console.log(total_sum);
}

solution_part_1();
solution_part_2()