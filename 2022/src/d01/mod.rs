extern crate core;

use itertools::Itertools;
use std::cmp::Reverse;

pub fn fn1_imperative(input: &str) -> i64 {
    let mut max = 0;

    for group in input.replace("\r\n", "\n").split("\n\n") {
        let mut sum = 0;
        for line in group.lines() {
            let value = line.parse::<i64>().unwrap();
            sum += value;
        }
        if sum > max {
            max = sum;
        }
    }

    max
}

pub fn fn1_partially_iterators(input: &str) -> i64 {
    // Works well but could potentially have large memory usage for big data sets.
    let lines = input
        .lines()
        .map(|v| v.parse::<i64>().ok())
        .collect::<Vec<_>>();

    let max = lines
        .split(|line| line.is_none())
        .map(|group| group.iter().map(|v| v.unwrap()).sum::<i64>())
        .max();

    max.unwrap()
}

pub fn fn1_fully_iterators_batching(input: &str) -> i64 {
    // Seems like the best approach.
    // Bounded memory usage for big data sets (moot at this point as baked-in to executable).
    input
        .lines()
        .map(|v| v.parse::<i64>().ok())
        .batching(|it| {
            let mut sum = None;
            while let Some(Some(v)) = it.next() {
                sum = Some(sum.unwrap_or(0) + v)
            }
            sum
        })
        .max()
        .unwrap_or_default()
}

pub fn fn1_fully_iterators_coalesce(input: &str) -> i64 {
    // Possibly impractical.
    input
        .lines()
        .map(|v| v.parse::<i64>().ok())
        .coalesce(|a, b| match (a, b) {
            (None, None) => Ok(None),
            (None, Some(b)) => Ok(Some(b)),
            (Some(a), Some(b)) => Ok(Some(a + b)),
            (Some(a), None) => Err((Some(a), None)),
        })
        .max()
        .flatten()
        .unwrap_or_default()
}

pub fn fn2_fully_iterators_batching(input: &str) -> i64 {
    input
        .lines()
        .map(|v| v.parse::<i64>().ok())
        .batching(|it| {
            let mut sum = None;
            while let Some(Some(v)) = it.next() {
                sum = Some(sum.unwrap_or(0) + v)
            }
            sum
        })
        .sorted_by_key(|&v| i64::MAX - v)
        .take(3)
        .sum::<i64>()
}

pub fn fn2_terse(input: &str) -> i64 {
    // I think I like this one best.
    input
        .lines()
        .map(|v| v.parse::<i64>().ok())
        .batching(|it| it.map_while(|x| x).sum1::<i64>())
        .map(Reverse)
        .k_smallest(3)
        .map(|x| x.0)
        .sum::<i64>()
}

#[cfg(test)]
mod tests {
    use super::*;
    use svutils::scaffold_test;

    const YEAR: i16 = 2022;
    const DAY: i16 = 1;

    #[test]
    fn test_fn1_example() {
        scaffold_test(
            YEAR,
            DAY,
            "example.txt",
            "example-spec.1.txt",
            fn1_imperative,
        );
        scaffold_test(
            YEAR,
            DAY,
            "example.txt",
            "example-spec.1.txt",
            fn1_partially_iterators,
        );
        scaffold_test(
            YEAR,
            DAY,
            "example.txt",
            "example-spec.1.txt",
            fn1_fully_iterators_batching,
        );
        scaffold_test(
            YEAR,
            DAY,
            "example.txt",
            "example-spec.1.txt",
            fn1_fully_iterators_coalesce,
        );
    }

    #[test]
    fn test_fn1_input() {
        scaffold_test(YEAR, DAY, "input.txt", "input-spec.1.txt", fn1_imperative);
        scaffold_test(
            YEAR,
            DAY,
            "input.txt",
            "input-spec.1.txt",
            fn1_partially_iterators,
        );
        scaffold_test(
            YEAR,
            DAY,
            "input.txt",
            "input-spec.1.txt",
            fn1_fully_iterators_batching,
        );
        scaffold_test(
            YEAR,
            DAY,
            "input.txt",
            "input-spec.1.txt",
            fn1_fully_iterators_coalesce,
        );
    }

    #[test]
    fn test_fn2_example() {
        scaffold_test(
            YEAR,
            DAY,
            "example.txt",
            "example-spec.2.txt",
            fn2_fully_iterators_batching,
        );
        scaffold_test(YEAR, DAY, "example.txt", "example-spec.2.txt", fn2_terse);
    }

    #[test]
    fn test_fn2_input() {
        scaffold_test(
            YEAR,
            DAY,
            "input.txt",
            "input-spec.2.txt",
            fn2_fully_iterators_batching,
        );
        scaffold_test(YEAR, DAY, "input.txt", "input-spec.2.txt", fn2_terse);
    }
}
