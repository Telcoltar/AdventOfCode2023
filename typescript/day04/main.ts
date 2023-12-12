import fs from "fs"

function parseNumberList(numbers: string): Set<number> {
  return new Set(numbers.trim().split(/\s+/).map((number: string) => parseInt(number)))
}
function loadData(filename: string): [Set<number>, Set<number>][] {
  let data = fs.readFileSync(filename)
  let lines = data.toString().split('\n')
  let cards: [Set<number>, Set<number>][] = []
  for (let line of lines) {
    let [_, numbers] = line.split(":")
    let [winningNumbers, myNumbers] = numbers.split("|").map(parseNumberList)
    cards.push([winningNumbers, myNumbers])
  }
  return cards
}

function solutionPart1() {
  let cards = loadData("input.txt")
  let total = 0
  for (let [winningNumbers, myNumbers] of cards) {
    let matches = [...winningNumbers].filter(x => myNumbers.has(x)).length
    if (matches > 0) {
      total += Math.pow(2, matches - 1)
    }
  }
  console.log(total)
}

function solutionPart2() {
  let cards = loadData("input.txt")
  let pile: number[] = Array.from<number>({ length: cards.length }).fill(1)
  for (let [index, [winningNumbers, myNumbers]] of cards.entries()) {
    let matches = [...winningNumbers].filter(x => myNumbers.has(x)).length
    for (let i = 0; i < matches; i++) {
      pile[i + index + 1] += pile[index]
    }
  }
  console.log(pile.reduce((a, b) => a + b, 0))
}

solutionPart1()
solutionPart2()