extern crate core;

use itertools::Itertools;

fn hash_fn(s: &str) -> i64 {
    s.as_bytes()
        .iter()
        .fold(0, |hash, &c| (hash + (c as u32)) * 17 % 256) as i64
}

#[derive(Debug)]
struct FocusHashmap<'a> {
    data: Vec<Vec<(&'a str, i64)>>,
}

impl<'a> FocusHashmap<'a> {
    fn new() -> Self {
        Self {
            data: vec![vec![]; 256],
        }
    }

    fn insert(&mut self, key: &'a str, value: i64) {
        let key_hash = hash_fn(key) as usize;
        if let Some(bucket) = self.data[key_hash].iter_mut().find(|(k, _)| *k == key) {
            bucket.1 = value;
        } else {
            self.data[key_hash].push((key, value));
        }
    }

    fn remove(&mut self, key: &'a str) {
        let key_hash = hash_fn(key) as usize;
        if let Some((index, _)) = self.data[key_hash]
            .iter()
            .find_position(|&&(k, _)| k == key)
        {
            self.data[key_hash].remove(index);
        }
    }

    fn total_power(&self) -> i64 {
        self.data
            .iter()
            .enumerate()
            .map(|(index, lens_box)| {
                (index as i64 + 1)
                    * lens_box
                        .iter()
                        .enumerate()
                        .map(|(slot, &(_, value))| (slot as i64 + 1) * value)
                        .sum::<i64>()
            })
            .sum::<i64>()
    }
}

pub fn fn1(input: &str) -> i64 {
    input.trim_end_matches('\n').split(',').map(hash_fn).sum()
}

pub fn fn2(input: &str) -> i64 {
    input
        .trim_end_matches('\n')
        .split(',')
        .fold(FocusHashmap::new(), |mut hashmap, s| {
            if s.contains('-') {
                let key = s.trim_end_matches('-');
                hashmap.remove(key);
            } else {
                // Do the elves not know to press shift?
                let (key, value) = s.split_once('=').unwrap();
                hashmap.insert(key, value.parse::<i64>().unwrap_or_default());
            }
            hashmap
        })
        .total_power() as i64
}

#[cfg(test)]
mod tests {
    use super::*;
    use svutils::scaffold_test;

    const YEAR: i16 = 2023;
    const DAY: i16 = 15;

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
