package solve2016

import "testing"

func referenceDragonExpansion(seed string, n int) string {
	for len(seed) < n {
		tail := ""
		for _, bit := range seed {
			flipped := "0"
			if bit == '0' {
				flipped = "1"
			}
			tail = flipped + tail
		}
		seed += "0" + tail
	}
	return seed[:n]
}

func TestDay16ExpansionMatchesReference(t *testing.T) {
	for length := 0; length <= 5; length++ {
		for bits := 0; bits < 1<<length; bits++ {
			seed := make([]byte, length)
			for i := range seed {
				seed[i] = '0' + byte((bits>>i)&1)
			}
			for target := 0; target <= 96; target++ {
				got := string((Day16{}).expand(string(seed), target))
				want := referenceDragonExpansion(string(seed), target)
				if got != want {
					t.Fatalf("seed %q target %d: got %q want %q", seed, target, got, want)
				}
			}
		}
	}
}

func referenceDragonChecksum(data string) string {
	for len(data) > 0 && len(data)%2 == 0 {
		next := ""
		for i := 0; i < len(data); i += 2 {
			if data[i] == data[i+1] {
				next += "1"
			} else {
				next += "0"
			}
		}
		data = next
	}
	return data
}

func TestDay16ChecksumMatchesReference(t *testing.T) {
	for length := 0; length <= 12; length++ {
		for bits := 0; bits < 1<<length; bits++ {
			data := make([]byte, length)
			for i := range data {
				data[i] = '0' + byte((bits>>i)&1)
			}
			want := referenceDragonChecksum(string(data))
			got := (Day16{}).checksum(data)
			if got != want {
				t.Fatalf("length %d bits %d: got %q want %q", length, bits, got, want)
			}
		}
	}
}

func TestDay16ExpansionAndChecksum(t *testing.T) {
	for _, seed := range []string{"0", "1", "110", "01101", "101001"} {
		for target := 1; target <= 128; target++ {
			want := referenceDragonChecksum(referenceDragonExpansion(seed, target))
			got := (Day16{}).checksum((Day16{}).expand(seed, target))
			if got != want {
				t.Fatalf("seed %q target %d: got %q want %q", seed, target, got, want)
			}
		}
	}
}

func TestDay16RejectsInvalidSeeds(t *testing.T) {
	for _, input := range []string{"", " \t\r\n", "102", "10x1", "10 1", "10\n1", "０１", "10\x00"} {
		for part, solve := range []func(string) (string, error){(Day16{}).Part1, (Day16{}).Part2} {
			if _, err := solve(input); err == nil {
				t.Errorf("part %d accepted %q", part+1, input)
			}
		}
	}
}

func TestDay16SolveTrimsAndValidatesWholeSeed(t *testing.T) {
	d := Day16{}
	for _, seed := range []string{"0", "1", "110", "01101"} {
		for _, size := range []int{1, 2, 9, 24, 272} {
			want := referenceDragonChecksum(referenceDragonExpansion(seed, size))
			got, err := d.solve(" \t"+seed+"\r\n", size)
			if err != nil || got != want {
				t.Fatalf("seed %q size %d: got %q, %v; want %q", seed, size, got, err, want)
			}
		}
	}
	if _, err := d.solve("110x", 2); err == nil {
		t.Fatal("accepted invalid character beyond truncation")
	}
}
