import * as fs from "fs"

enum Direction {
    UP = 1,
    DOWN = 2,
    LEFT = 4,
    RIGHT = 8
}

function getOppositeDirection(direction: Direction): Direction {
    switch (direction) {
        case Direction.UP:
            return Direction.DOWN
        case Direction.DOWN:
            return Direction.UP
        case Direction.LEFT:
            return Direction.RIGHT
        case Direction.RIGHT:
            return Direction.LEFT
    }
}

type Point = {
    x: number,
    y: number
}

function movePoint(point: Point, direction: Direction) {
    switch (direction) {
        case Direction.UP:
            point.y -= 1
            break
        case Direction.DOWN:
            point.y += 1
            break
        case Direction.LEFT:
            point.x -= 1
            break
        case Direction.RIGHT:
            point.x += 1
            break
    }
}

type Status = {
    direction: Direction,
    position: Point
}

function moveStatus(status: Status) {
    movePoint(status.position, status.direction)
}

function moveStatusBack(status: Status) {
    movePoint(status.position, getOppositeDirection(status.direction))
}

function loadData(filename: string): string[][] {
    return fs.readFileSync(filename, "utf-8").split("\n")
        .map(line => line.trim().split(""))
}

function padGrid(grid: string[][]): string[][] {
    const newGrid = grid.map(line => ["#"].concat(line).concat(["#"]))
    const emptyLine = new Array(newGrid[0].length).fill("#")
    return [emptyLine].concat(newGrid).concat([emptyLine])

}

function processForwardSlash(status: Status) {
    switch (status.direction) {
        case Direction.UP:
            status.direction = Direction.RIGHT
            break
        case Direction.DOWN:
            status.direction = Direction.LEFT
            break
        case Direction.LEFT:
            status.direction = Direction.DOWN
            break
        case Direction.RIGHT:
            status.direction = Direction.UP
            break
    }
}

function processBackSlash(status: Status) {
    switch (status.direction) {
        case Direction.UP:
            status.direction = Direction.LEFT
            break
        case Direction.DOWN:
            status.direction = Direction.RIGHT
            break
        case Direction.LEFT:
            status.direction = Direction.UP
            break
        case Direction.RIGHT:
            status.direction = Direction.DOWN
            break
    }
}

function processLeftRight(status: Status): Status[] {
    switch (status.direction) {
        case Direction.UP:
        case Direction.DOWN:
            const left = {...status, position: {...status.position}}
            const right = {...status, position: {...status.position}}
            left.direction = Direction.LEFT
            right.direction = Direction.RIGHT
            return [left, right]
        case Direction.LEFT:
        case Direction.RIGHT:
            return [status]
    }
}

function processUpDown(status: Status): Status[] {
    switch (status.direction) {
        case Direction.UP:
        case Direction.DOWN:
            return [status]
        case Direction.LEFT:
        case Direction.RIGHT:
            const up = {...status, position: {...status.position}}
            const down = {...status, position: {...status.position}}
            up.direction = Direction.UP
            down.direction = Direction.DOWN
            return [up, down]
    }
}

function countGrid(grid: boolean[][]): number {
    return grid.map(line => line.filter(value => value).length).reduce((a, b) => a + b, 0)
}

function calculateEnergizedTiles(grid: string[][], start: Status): number {
    const visited = Array.from(
        { length: grid.length },
        () => new Uint8Array(grid.length)
    )
    const count: boolean[][] = Array.from({ length: grid.length }, () => new Array(grid.length).fill(false))
    const queue: Status[] = [start]
    while (queue.length > 0) {
        const status = queue.shift()!
        if (!(visited[status.position.y][status.position.x] & status.direction)) {
            visited[status.position.y][status.position.x] = visited[status.position.y][status.position.x] | status.direction
            count[status.position.y][status.position.x] = true
            moveStatus(status)
            const tile = grid[status.position.y][status.position.x]
            switch (tile) {
                case "/":
                    processForwardSlash(status)
                    queue.push(status)
                    break
                case "\\":
                    processBackSlash(status)
                    queue.push(status)
                    break
                case "-":
                    processLeftRight(status).forEach(s => queue.push(s))
                    break
                case "|":
                    processUpDown(status).forEach(s => queue.push(s))
                    break
                case ".":
                    while (grid[status.position.y][status.position.x] === ".") {
                        count[status.position.y][status.position.x] = true
                        moveStatus(status)
                    }
                    moveStatusBack(status)
                    queue.push(status)
                    break
            }
        }
    }
    return countGrid(count) - 1
}

function solutionPart1() {
    let grid = loadData("input.txt")
    grid = padGrid(grid)
    const startTime = new Date().getTime()
    console.log(calculateEnergizedTiles(grid, {position: {x: 0, y: 1}, direction: Direction.RIGHT}))
    const endTime = new Date().getTime()
    console.log(`Time: ${endTime - startTime}ms`)
}

function solutionPart2() {
    let grid = loadData("input.txt")
    grid = padGrid(grid)
    const startTime = new Date().getTime()
    let energizedTiles: number[] = []
    for (let i = 1; i < grid.length - 1; i++) {
        energizedTiles.push(
            calculateEnergizedTiles(grid, {position: {x: 0, y: i}, direction: Direction.RIGHT})
        )
        energizedTiles.push(
            calculateEnergizedTiles(grid, {position: {x: grid.length - 1, y: i}, direction: Direction.LEFT})
        )
        energizedTiles.push(
            calculateEnergizedTiles(grid, {position: {x: i, y: 0}, direction: Direction.DOWN})
        )
        energizedTiles.push(
            calculateEnergizedTiles(grid, {position: {x: i, y: grid.length - 1}, direction: Direction.UP})
        )
    }
    console.log(Math.max(...energizedTiles))
    const endTime = new Date().getTime()
    console.log(`Time: ${endTime - startTime}ms`)
}

solutionPart1()
solutionPart2()
