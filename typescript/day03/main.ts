import fs from "fs"

type Point = [number, number]
type Symbol = [string, Point]
function isDigit(char: string): boolean {
  const charCode = char.charCodeAt(0)
  return charCode >= 48 && charCode <= 57 // Unicode values for digits 0-9 are 48-57
}
function load_board(filename: string) {
  let data = fs.readFileSync(filename)
  let lines = data.toString().split('\n')
  let board: string[] = []
  let width = lines[0].length + 2
  board.push(Array(width).fill('.').join(""))
  lines.forEach((line: string) => {
    board.push('.' + line + '.')
  })
  board.push(Array(width).fill('.').join(""))
  let numbers: [number, Symbol[]][] = []
  let y = 0
  let start: Point = [0, 0]
  for (let line of board) {
    let number_in_progress = false
    let current_number = ""
    let x = 0
    for (let char of line) {
      if (isDigit(char)) {
        if (number_in_progress) {
          current_number += char
        } else {
          number_in_progress = true
          current_number = char
          start = [x, y]
        }
      } else {
        if (number_in_progress) {
          number_in_progress = false
          numbers.push(
            [parseInt(current_number), get_symbols_in_surrounding(board, start, current_number.length)]
          )
        }
      }
      x++
    }
    y++
  }
  return numbers
}

function get_symbols_in_surrounding(board: string[], start: Point, length: number) {
  let symbols: Symbol[] = []
  let [x, y] = start
  for (let i = -1; i < length + 1; i++) {
    if (board[y + 1][x + i] != '.' && !isDigit(board[y + 1][x + i])) {
      symbols.push([board[y + 1][x + i], [x + i, y + 1]])
    }
    if (board[y - 1][x + i] != '.' && !isDigit(board[y - 1][x + i])) {
      symbols.push([board[y - 1][x + i], [x + i, y - 1]])
    }
  }
  if (board[y][x - 1] != '.' && !isDigit(board[y][x -1])) {
    symbols.push([board[y][x - 1], [x - 1, y]])
  }
  if (board[y][x + length] != '.' && !isDigit(board[y][x + length])) {
    symbols.push([board[y][x + length], [x + length, y]])
  }
  return symbols
}

function solution_part_1() {
  let numbers = load_board("input.txt")
  let part_number_sum = 0
  for (let [number, symbols] of numbers) {
    if (symbols.length > 0) {
      part_number_sum += number
    }
  }
  console.log(part_number_sum)
}

function solution_part_2() {
  let numbers = load_board("input.txt")
  let stars: Map<string, number[]> = new Map()
  for (let [number, symbols] of numbers) {
    for (let [symbol, point] of symbols) {
      if (symbol == "*") {
        let number_list = stars.get(point[0] + "." + point[1]) ?? []
        number_list.push(number)
        stars.set(point[0] + "." + point[1], number_list)
        break
      }
    }
  }
  let gear_ratio_sum = 0
  for (let [_, numbers] of stars) {
    if (numbers.length == 2) {
      gear_ratio_sum += numbers[0] * numbers[1]
    }
  }
  console.log(gear_ratio_sum)
}

solution_part_1()
solution_part_2()