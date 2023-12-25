extern crate core;

use csv::ReaderBuilder;
use itertools::Itertools;
use std::io::Cursor;
use std::collections::HashMap;

#[derive(Debug)]
struct CubeDraw {
    red: i64,
    green: i64,
    blue: i64,
}

impl CubeDraw {
    fn new(value: String) -> Self {
        let mut counts = HashMap::new();
        for part in value.split(", ") {
            let mut split = part.split_whitespace();
            let count: i64 = split.next().unwrap_or_default().parse().unwrap_or_default();
            let color = split.next().unwrap_or_default();
            *counts.entry(color).or_insert(0) += count;
        }

        CubeDraw {
            red: *counts.get("red").unwrap_or(&0),
            green: *counts.get("green").unwrap_or(&0),
            blue: *counts.get("blue").unwrap_or(&0),
         }
    }
}

fn parse_game(games: Vec<String>) -> Vec<CubeDraw> {
    games
        .into_iter()
        .map(CubeDraw::new)
        .collect_vec()
}

fn parse_games(input: &str) -> HashMap<i64, Vec<CubeDraw>> {
    let mut reader = ReaderBuilder::new()
        .has_headers(false)
        .delimiter(b':')
        .from_reader(Cursor::new(input));

    let mut map: HashMap<i64, Vec<CubeDraw>> = HashMap::new();

    for result in reader.records() {
        let record = result.unwrap_or_default();

        let game_key = record
            .get(0)
            .unwrap_or_default()
            .strip_prefix("Game ")
            .unwrap_or_default()
            .parse::<i64>()
            .unwrap_or_default();

        let game_strs = record
            .get(1)
            .unwrap_or_default()
            .trim_start()
            .to_string()
            .split("; ")
            .map(|s| s.to_string())
            .collect_vec();

        let games = parse_game(game_strs);

        map.insert(game_key, games);
    };

    map
}

pub fn fn1(input: &str) -> i64 {
    let all_games = parse_games(input);

    let valid_games = all_games.iter().filter(|g| {
        let game_has_invalid_draws = g.1.iter().any(|draw| {
            if draw.blue > 14 {
                return true
            }

            if draw.green > 13 {
                return true
            }

            if draw.red > 12 {
                return true
            }

            return false
        });
        
        return !game_has_invalid_draws
    });

    valid_games.map(|g| g.0).sum()
}

pub fn fn2(input: &str) -> i64 {
    let all_games = parse_games(input);

    all_games.iter().map(|g| {
        let mut biggest_draw = CubeDraw{
            red: 0,
            green: 0,
            blue: 0,
        };

        for x in g.1.iter() {
            biggest_draw = CubeDraw{
                red: if x.red > biggest_draw.red { x.red } else { biggest_draw.red },
                green: if x.green > biggest_draw.green { x.green } else { biggest_draw.green },
                blue: if x.blue > biggest_draw.blue { x.blue } else { biggest_draw.blue },
            }
        }

        biggest_draw
    }).map(|biggest_draw| {
        biggest_draw.red * biggest_draw.green * biggest_draw.blue
    })
    .sum()
}

#[cfg(test)]
mod tests {
    use super::*;
    use svutils::load_spec;

    #[test]
    fn test_fn1_example() {
        assert_eq!(fn1(include_str!("../../../input/2023/d02/example.txt")), load_spec(include_str!("../../../input/2023/d02/example-spec.1.txt")));
    }

    #[test]
    fn test_fn1_input() {
        assert_eq!(fn1(include_str!("../../../input/2023/d02/input.txt")), load_spec(include_str!("../../../input/2023/d02/input-spec.1.txt")));
    }

    #[test]
    fn test_fn2_example() {
        assert_eq!(fn2(include_str!("../../../input/2023/d02/example.txt")), load_spec(include_str!("../../../input/2023/d02/example-spec.2.txt")));
    }

    #[test]
    fn test_fn2_input() {
        assert_eq!(fn2(include_str!("../../../input/2023/d02/input.txt")), load_spec(include_str!("../../../input/2023/d02/input-spec.2.txt")));
    }
}