package solve2024

import "testing"

func TestDay10CountsPeaksAndRoutesSeparately(t *testing.T) {
	for _, tc := range []struct{ name, input, peaks, routes string }{
		{"single route", "0123456789", "1", "1"},
		{"merged routes", "0123456788\n1234567898", "1", "9"},
		{"two destinations", "0123456789\n1234567899", "2", "10"},
		{"shared destination across starts", "0123456789\n1011111111", "2", "3"},
		{"no complete route", "0123456788", "0", "0"},
		{"no start", "1234567899", "0", "0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for part, check := range []struct {
				solve func(string) (string, error)
				want  string
			}{
				{(Day10{}).Part1, tc.peaks}, {(Day10{}).Part2, tc.routes},
			} {
				got, err := check.solve(tc.input)
				if err != nil || got != check.want {
					t.Fatalf("part %d got %q, %v; want %q", part+1, got, err, check.want)
				}
			}
		})
	}
}

func TestDay10RejectsInvalidMaps(t *testing.T) {
	for _, input := range []string{
		"", "\n", "0", "012345678", "012\n345\n678",
		"0123456789\n012345678", "0123456789\n\n0123456789",
		"0123456789\n\n", "\n0123456789", "01234 6789",
		"012345678x", "012345678９", "0123456789\r", "012345678\x00",
	} {
		for part, solve := range []func(string) (string, error){(Day10{}).Part1, (Day10{}).Part2} {
			if _, err := solve(input); err == nil {
				t.Errorf("part %d accepted %q", part+1, input)
			}
		}
	}
}

func TestDay10AcceptsLineEndings(t *testing.T) {
	for _, input := range []string{
		"0123456789", "0123456789\n", "0123456789\r\n",
		"01234\n98765", "01234\n98765\n", "01234\r\n98765\r\n",
	} {
		for part, solve := range []func(string) (string, error){(Day10{}).Part1, (Day10{}).Part2} {
			got, err := solve(input)
			if err != nil || got != "1" {
				t.Errorf("part %d input %q: got %q, %v; want 1", part+1, input, got, err)
			}
		}
	}
}
