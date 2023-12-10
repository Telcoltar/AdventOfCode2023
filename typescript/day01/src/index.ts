import {readFileSync} from 'fs'
import * as console from "console";

// english digit word to digit conversion
const NUMBERS: Record<string, string> = {
  "one": "1",
  "two": "2",
  "three": "3",
  "four": "4",
  "five": "5",
  "six": "6",
  "seven": "7",
  "eight": "8",
  "nine": "9"
}
function isDigit(char: string): boolean {
    const charCode = char.charCodeAt(0);
    return charCode >= 48 && charCode <= 57; // Unicode values for digits 0-9 are 48-57
}
function solution_part_1() {
  // reading file
  let data = readFileSync("input.txt")
  // line_numbers array
  let line_numbers: number[] = [];
  // loop over lines
  data.toString().split('\n').forEach((line: string) => {
    let line_number = ""
    // loop over chars
    for (let char of line) {
      // check if char is a digit
      if (isDigit(char)) {
        line_number += char
        break
      }
    }
    // loop over chars in reverse
    for (let i = line.length - 1; i >= 0; i--) {
        let char = line[i];
        if (isDigit(char)) {
            line_number += char;
            break;
        }
    }
    // parse line_number and append to line_numbers
    line_numbers.push(parseInt(line_number));
  })
  const sum = line_numbers.reduce((a, b) => a + b, 0);
  console.log(sum);
}

function extractFirstDigit(line: string): string {
  let constructedNumber = "";
  const firstDigit = line.match(/(\d|one|two|three|four|five|six|seven|eight|nine)/);
  if (firstDigit) {
    constructedNumber += firstDigit[0].length > 1 ? NUMBERS[firstDigit[0]] : firstDigit[0];
  }
  return constructedNumber;
}

function extractLastDigit(line: string): string {
  let constructedNumber = "";
  const reversedLine = line.split('').reverse().join('');
  const lastDigitReversed = reversedLine.match(/(\d|eno|owt|eerht|ruof|evif|xis|neves|thgie|enin)/);
  if (lastDigitReversed) {
    const lastDigit = lastDigitReversed[0].split('').reverse().join('');
    constructedNumber += lastDigit.length > 1 ? NUMBERS[lastDigit] : lastDigit;
  }
  return constructedNumber;
}

function processLine(line: string): number {
  let constructedNumber = extractFirstDigit(line);
  constructedNumber += extractLastDigit(line);
  return parseInt(constructedNumber);
}

function sumNumbers(lineNumbers: number[]): number {
  return lineNumbers.reduce((a, b) => a + b, 0);
}

function solution_part_2() {
  const fileData = readFileSync("input.txt");
  const lineNumbers: number[] = fileData.toString().split('\n').map(processLine);
  const sum = sumNumbers(lineNumbers);
  console.log(sum);
}

solution_part_2()