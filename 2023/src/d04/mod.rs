extern crate core;

struct Card {
    count: u32,
    winning: Vec<u32>,
    held: Vec<u32>,
}

fn parse_card(input: &str) -> Card {
    let (_, card_part) = input.split_once(": ").unwrap();
    let (winning_part, held_part) = card_part.split_once("|").unwrap();
    let winning = winning_part
        .split_whitespace()
        .map(|s| s.parse::<u32>().unwrap())
        .collect();
    let held = held_part
        .split_whitespace()
        .map(|s| s.parse::<u32>().unwrap())
        .collect();
    let count = 1;
    Card {
        winning,
        held,
        count,
    }
}

fn parse_all(input: &str) -> Vec<Card> {
    input.lines().map(parse_card).collect()
}

pub fn fn1(input: &str) -> i64 {
    let cards = parse_all(&input);
    let mut total = 0;
    for card in cards.iter() {
        let numbers_won = card
            .held
            .iter()
            .filter(|num| card.winning.contains(num))
            .count();
        if numbers_won > 0 {
            total += 1 << (numbers_won - 1);
        }
    }
    total
}

pub fn fn2(input: &str) -> i64 {
    let mut cards = parse_all(&input);
    let len = cards.len();
    for idx in 0..len {
        let card = &cards[idx];
        let numbers_won = card
            .held
            .iter()
            .filter(|num| card.winning.contains(num))
            .count();
        let n = card.count;
        for i in 0..numbers_won {
            cards[idx + i + 1].count += n
        }
    }
    cards.iter().map(|card| card.count as i64).sum::<i64>()
}

#[cfg(test)]
mod tests {
    use super::*;
    use svutils::scaffold_test;

    const YEAR: i16 = 2023;
    const DAY: i16 = 4;

    #[test]
    fn test_fn1_example() {
        scaffold_test(YEAR, DAY, "example.txt", "example-spec.1.txt", fn1);
    }

    #[test]
    fn test_fn1_input() {
        scaffold_test(YEAR, DAY, "input.txt", "input-spec.1.txt", fn1);
    }

    #[test]
    fn test_fn2_example() {
        scaffold_test(YEAR, DAY, "example.txt", "example-spec.2.txt", fn2);
    }

    #[test]
    fn test_fn2_input() {
        scaffold_test(YEAR, DAY, "input.txt", "input-spec.2.txt", fn2);
    }
}
