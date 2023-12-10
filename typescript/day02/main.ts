import fs from 'fs'

const CUBES_MAX: Map<string, number> = new Map([
    ["red", 12],
    ["green", 13],
    ["blue", 14]
])
function load_game_data(filename: string): Map<number, Map<string, number>> {
    let game_data = new Map<number, Map<string, number>>();
    let data = fs.readFileSync(filename)
    data.toString().split('\n').forEach((line: string) => {
        let [game, draws] = line.split(':')
        let [_, game_index_str] = game.split(' ')
        let game_index = parseInt(game_index_str)
        let cubes_max: Map<string, number> = new Map([
            ["red", 0],
            ["green", 0],
            ["blue", 0]
        ])
        for (let draw of draws.split(';')) {
            for (let cubes of draw.split(',')) {
                let [cube_count_str, cube_color] = cubes.trim().split(' ')
                let cube_count = parseInt(cube_count_str)
                cubes_max.set(cube_color, Math.max(cubes_max.get(cube_color) ?? 0, cube_count))
            }
        }
        game_data.set(game_index, cubes_max)
    })
    return game_data
}

function solution_part_1() {
    let possible_games_sum = 0
    for (let [game_index, game] of load_game_data("input.txt")) {
        if (
          game.get("red")! <= CUBES_MAX.get("red")! &&
            game.get("green")! <= CUBES_MAX.get("green")! &&
            game.get("blue")! <= CUBES_MAX.get("blue")!
        ) {
            possible_games_sum += game_index
        }
    }
    console.log(possible_games_sum)
}

function solution_part_2() {
    let power_sum = 0
    for (let [_, game] of load_game_data("input.txt")) {
        power_sum += game.get("red")! * game.get("green")! * game.get("blue")!
    }
    console.log(power_sum)
}

solution_part_1()
solution_part_2()