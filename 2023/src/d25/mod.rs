extern crate core;

pub fn fn1(input: &str) -> i64 {
    // Build a list of connections
    // Initial solve in Python/networkx. Will try to implement the solution in Rust soon.

    todo!()
}

pub fn fn2(_input: &str) -> i64 {
    todo!()
}

#[cfg(test)]
mod tests {
    use super::*;
    use svutils::load_spec;

    #[test]
    fn test_fn1_example() {
        assert_eq!(fn1(include_str!("../../../input/2023/d25/example.txt")), load_spec(include_str!("../../../input/2023/d25/example-spec.1.txt")));
    }

    #[test]
    fn test_fn1_input() {
        assert_eq!(fn1(include_str!("../../../input/2023/d25/input.txt")), load_spec(include_str!("../../../input/2023/d25/input-spec.1.txt")));
    }

    #[test]
    fn test_fn2_example() {
        assert_eq!(fn2(include_str!("../../../input/2023/d25/example.txt")), load_spec(include_str!("../../../input/2023/d25/example-spec.2.txt")));
    }

    #[test]
    fn test_fn2_input() {
        assert_eq!(fn2(include_str!("../../../input/2023/d25/input.txt")), load_spec(include_str!("../../../input/2023/d25/input-spec.2.txt")));
    }

}