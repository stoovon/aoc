extern crate core;

pub fn fn1(input: &str) -> i64 {
    input
        .chars()
        .map(|c| match c {
            '(' => 1,
            ')' => -1,
            _ => 0,
        })
        .fold(0, |sum, i| sum + i)
}

pub fn fn2(input: &str) -> i64 {
    let mut n = 0;
    for (i, c) in input.chars().enumerate() {
        match c {
            '(' => n += 1,
            ')' => n -= 1,
            _ => (),
        }
        if n == -1 {
            return i as i64 + 1;
        }
    }

    0
}

#[cfg(test)]
mod tests {
    use super::*;
    use svutils::scaffold_test;

    const YEAR: i16 = 2015;
    const DAY: i16 = 1;

    #[test]
    fn test_fn1_unit() {
        assert_eq!(fn1("(())"), 0);
        assert_eq!(fn1("()()"), 0);
        assert_eq!(fn1("((("), 3);
        assert_eq!(fn1("(()(()("), 3);
        assert_eq!(fn1("))((((("), 3);
        assert_eq!(fn1("())"), -1);
        assert_eq!(fn1("))("), -1);
        assert_eq!(fn1(")))"), -3);
        assert_eq!(fn1(")())())"), -3);
    }

    #[test]
    fn test_fn1_input() {
        scaffold_test(YEAR, DAY, "input.txt", "input-spec.1.txt", fn1);
    }

    #[test]
    fn test_fn2_unit() {
        assert_eq!(fn2(")"), 1);
        assert_eq!(fn2("()())"), 5);
    }

    #[test]
    fn test_fn2_input() {
        scaffold_test(YEAR, DAY, "input.txt", "input-spec.2.txt", fn2);
    }
}
